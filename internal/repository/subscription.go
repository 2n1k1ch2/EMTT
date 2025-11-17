package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var errNotFound = errors.New("Subscription not found")

func NewSubScriptRepo(db *pgxpool.Pool, ctx context.Context) *SubScriptRepo {
	return &SubScriptRepo{db: db, ctx: ctx}
}

func (r *SubScriptRepo) GetById(id pgtype.UUID) (*Subscription, error) {
	var sub Subscription
	row := r.db.QueryRow(r.ctx,
		`SELECT id, user_id, service_name, price, start_date, finish_date
		 FROM subscription WHERE id = $1`, id)

	err := row.Scan(
		&sub.ID,
		&sub.UserID,
		&sub.ServiceName,
		&sub.Price,
		&sub.StartDate,
		&sub.FinishDate,
	)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *SubScriptRepo) Create(sub *Subscription) error {
	_, err := r.db.Exec(r.ctx,
		`INSERT INTO subscription (id, user_id, service_name, price, start_date, finish_date)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		sub.ID,
		sub.UserID,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.FinishDate,
	)

	return err
}

func (r *SubScriptRepo) Update(sub *Subscription) error {
	cmd, err := r.db.Exec(r.ctx,
		`UPDATE subscription
		 SET service_name = $1, price = $2, start_date = $3, finish_date = $4
		 WHERE id = $5`,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.FinishDate,
		sub.ID,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errNotFound
	}
	return nil
}

func (r *SubScriptRepo) Delete(id pgtype.UUID) error {
	cmd, err := r.db.Exec(r.ctx, `DELETE FROM subscription WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errNotFound
	}
	return nil
}

func (r *SubScriptRepo) List() ([]*Subscription, error) {
	rows, err := r.db.Query(r.ctx,
		`SELECT id, user_id, service_name, price, start_date, finish_date
		 FROM subscription`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*Subscription

	for rows.Next() {
		var sub Subscription
		err = rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.ServiceName,
			&sub.Price,
			&sub.StartDate,
			&sub.FinishDate,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, &sub)
	}

	return subs, nil
}

func (r *SubScriptRepo) TotalCost(userID string, service, from, to string) (uint32, error) {
	var count uint32

	rows, err := r.db.Query(
		r.ctx,
		`SELECT COALESCE(SUM(price), 0)
		 FROM subscription
		 WHERE ($1 = '' OR user_id = $1::uuid)
		   AND ($2 = '' OR service_name ILIKE '%' || $2 || '%')
		   AND ($3 = '' OR start_date >= $3::date)
		   AND ($4 = '' OR finish_date <= $4::date OR finish_date IS NULL)
		`,
		userID, service, from, to,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	return count, nil
}
