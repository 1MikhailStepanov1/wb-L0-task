package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

func NewConnectionPool(config *Config) (*pgxpool.Pool, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.User, config.Password, config.Host, config.Port, config.Database,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open postgres connection: %w", err)
	}
	return pool, nil
}

func Stop(p *pgxpool.Pool) {
	p.Close()
}
