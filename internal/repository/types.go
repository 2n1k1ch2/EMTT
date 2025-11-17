package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Subscription struct {
	ID          pgtype.UUID
	ServiceName string
	Price       int
	UserID      pgtype.UUID
	StartDate   pgtype.Timestamptz
	FinishDate  pgtype.Timestamptz
}

type SubscriptionRepository interface {
	GetById(ID pgtype.UUID) (*Subscription, error)
	Create(sub *Subscription) error
	Update(sub *Subscription) error
	Delete(id pgtype.UUID) error
	List() ([]*Subscription, error)
	TotalCost(userID string, service, from, to string) (uint32, error)
}

type SubScriptRepo struct {
	db  *pgxpool.Pool
	ctx context.Context
}
