package postgres

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	TrManager interface {
		Do(ctx context.Context, fn func(ctx context.Context) error) error
	}
)

type Repo struct {
	db        *pgxpool.Pool
	getter    *trmpgx.CtxGetter
	trManager TrManager
}

func NewRepo(db *pgxpool.Pool, trManager TrManager, c *trmpgx.CtxGetter) *Repo {
	return &Repo{
		db:        db,
		getter:    c,
		trManager: trManager,
	}
}
