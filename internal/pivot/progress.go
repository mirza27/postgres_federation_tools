package pivot

import (
	"context"
	"database/sql"

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
