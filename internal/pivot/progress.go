package pivot

import (
	"context"
	"database/sql"
	"db_migrate_server/internal/util"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type _ExecSplit struct {
	SQLText string
	Status  string
	SQLArgs string
}

type _ExecQueue struct {
	QueueID   uuid.UUID
	Entity    string
	SQLText   string
	SQLArgs   string
	Status    string
	ExecSplit []_ExecSplit
	LastError sql.NullString
}

func (r *Repo) GetExecutionQueueList(limit int, page int, searchSQLText string, searchSQLArg string, filterStatus []string, filterEntity []string) ([]_ExecQueue, error) {
	ctx := context.Background()

	query := `SELECT DISTINCT q.queue_id, q.entity, q.sql_text, q.sql_args, q.status, q.last_error FROM _exec_queue q`
	args := []interface{}{}
	argIndex := 1

	// If search is provided, add LEFT JOIN with _exec_split
	if searchSQLText != "" || searchSQLArg != "" {
		query += ` LEFT JOIN _exec_split s ON q.queue_id = s.queue_id`
	}

	query += ` WHERE 1=1`

	// Filter by search sql_text (search in both queue and split)
	if searchSQLText != "" {
		query += fmt.Sprintf(" AND (q.sql_text ILIKE $%d OR s.sql_text ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+searchSQLText+"%")
		argIndex++
	}

	// Filter by search value (search in both queue and split)
	if searchSQLArg != "" {
		query += fmt.Sprintf(" AND (q.sql_args::text ILIKE $%d OR s.sql_args::text ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+searchSQLArg+"%")
		argIndex++
	}

	// Filter by status
	if len(filterStatus) > 0 {
		placeholders := []string{}
		for _, status := range filterStatus {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, status)
			argIndex++
		}
		query += " AND q.status IN (" + strings.Join(placeholders, ",") + ")"
	}

	// Filter by entity
	if len(filterEntity) > 0 {
		placeholders := []string{}
		for _, entity := range filterEntity {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIndex))
			args = append(args, entity)
			argIndex++
		}
		query += " AND q.entity IN (" + strings.Join(placeholders, ",") + ")"
	}

	// Add pagination
	offset := (page - 1) * limit
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	util.Debug.Printf("GetExecutionQueueList: query=%s args=%v", query, args)

	// Execute query
	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []_ExecQueue
	for rows.Next() {
		var q _ExecQueue
		if scanErr := rows.Scan(&q.QueueID, &q.Entity, &q.SQLText, &q.SQLArgs, &q.Status, &q.LastError); scanErr != nil {
			continue
		}

		// Fetch splits for this queue
		srows, serr := r.DB.Query(ctx, `
			select sql_text, sql_args, status
			from _exec_split where queue_id=$1
			order by created_at asc`, q.QueueID)
		if serr == nil {
			for srows.Next() {
				var s _ExecSplit
				if scanErr := srows.Scan(&s.SQLText, &s.SQLArgs, &s.Status); scanErr != nil {
					continue
				}
				q.ExecSplit = append(q.ExecSplit, s)
			}
			srows.Close()
		}

		out = append(out, q)
	}

	return out, nil
}

func (r *Repo) GetLastUpdatedQueueList(limit int) ([]_ExecQueue, error) {
	ctx := context.Background()

	// First pick the most recently updated queue_ids (considering splits)
	idRows, err := r.DB.Query(ctx, `
		select q.queue_id
		from _exec_queue q
		left join _exec_split s on s.queue_id = q.queue_id
		group by q.queue_id, q.updated_at
		order by coalesce(max(s.updated_at), q.updated_at) desc
		limit $1`, limit)
	if err != nil {
		return nil, err
	}
	defer idRows.Close()

	var ids []uuid.UUID
	for idRows.Next() {
		var id uuid.UUID
		if scanErr := idRows.Scan(&id); scanErr != nil {
			continue
		}
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return []_ExecQueue{}, nil
	}

	var out []_ExecQueue
	for _, qid := range ids {
		var q _ExecQueue
		// fetch queue row
		if err := r.DB.QueryRow(ctx, `
			select queue_id, entity, sql_text, sql_args,status, last_error
			from _exec_queue where queue_id=$1`, qid).Scan(&q.QueueID, &q.Entity, &q.SQLText, &q.SQLArgs, &q.Status, &q.LastError); err != nil {
			// skip if cannot fetch
			continue
		}

		// fetch splits for this queue
		srows, serr := r.DB.Query(ctx, `
			select sql_text, sql_args, status
			from _exec_split where queue_id=$1
			order by created_at asc`, qid)
		if serr == nil {
			for srows.Next() {
				var s _ExecSplit
				if scanErr := srows.Scan(&s.SQLText, &s.SQLArgs, &s.Status); scanErr != nil {
					continue
				}
				q.ExecSplit = append(q.ExecSplit, s)
			}
			srows.Close()
		}

		out = append(out, q)
	}

	return out, nil
}

type ExecQueueSummary struct {
	Total   int `json:"total"`
	Ready   int `json:"ready"`
	Pending int `json:"pending"`
	Error   int `json:"error"`
	Done    int `json:"done"`
}

func (r *Repo) GetExecQueueSummary() (ExecQueueSummary, error) {

	ctx := context.Background()

	var total, ready, pending, errCnt, done int
	row := r.DB.QueryRow(ctx, `
		select
			count(*) as total,
			count(*) filter (where status='ready') as ready,
			count(*) filter (where status='pending') as pending,
			count(*) filter (where status='error') as error,
			count(*) filter (where status='done') as done
		from _exec_queue`)
	if scanErr := row.Scan(&total, &ready, &pending, &errCnt, &done); scanErr != nil {
		return ExecQueueSummary{}, scanErr
	}

	return ExecQueueSummary{
		Total:   total,
		Ready:   ready,
		Pending: pending,
		Error:   errCnt,
		Done:    done,
	}, nil
}
