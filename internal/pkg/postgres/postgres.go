package postgres

import (
	"context"
	"fmt"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func SetupPostgres(ctx context.Context, c *Config) (*pgxpool.Pool, *manager.Manager, *trmpgx.CtxGetter, error) {
	pgCfg, err := pgxpool.ParseConfig(fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.User, c.Password, c.Host, c.Port, c.Database,
	))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("pgxpool.ParseConfig failed")
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("pgxpool.NewWithConfig failed")
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf("postgres.Ping failed")
	}

	trManager := manager.Must(trmpgx.NewDefaultFactory(pool))
	ctxGetter := trmpgx.DefaultCtxGetter
	return pool, trManager, ctxGetter, nil
}
