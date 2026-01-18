package pivot

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type _ExecQueue struct {
	QueueID   uuid.UUID
	Entity    string
	SQLText   string
	Status    string
	LastError sql.NullString
}

func (r *Repo) GetLastUpdatedQueueList(limit int) ([]_ExecQueue, error) {

	ctx := context.Background()

	rows, err := r.DB.Query(ctx, `
		select q.queue_id, q.entity,
			   coalesce(s.sql_text, q.sql_text) as sql_text,
			   coalesce(s.status, q.status) as status,
			   coalesce(s.last_error, q.last_error) as last_error,
			   greatest(coalesce(s.updated_at, q.updated_at), q.updated_at) as last_upd
		from _exec_queue q
		left join _exec_split s on s.queue_id = q.queue_id
		order by last_upd desc
		limit $1`, limit)
	if err != nil {
		return []_ExecQueue{}, err
	}
	defer rows.Close()

	var out []_ExecQueue
	for rows.Next() {
		var (
			queueID uuid.UUID
			entity  string
			sqlText string
			status  string
			lastErr sql.NullString
			_dummy  interface{}
		)
		if err := rows.Scan(&queueID, &entity, &sqlText, &status, &lastErr, &_dummy); err != nil {
			continue
		}
		// if lastErr.Valid {
		// 	errMsg = lastErr.String
		// }
		out = append(out, _ExecQueue{
			QueueID:   queueID,
			Entity:    entity,
			SQLText:   sqlText,
			Status:    status,
			LastError: lastErr,
		})
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
