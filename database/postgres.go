package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/Walker088/gorealestate/config"
	e "github.com/Walker088/gorealestate/error"
)

const (
	currentPackage = "github.com/Walker088/gorealestate/database"

	DbConnStringParsingError = "DB00001"
	DbTxPoolCreatingError    = "DB00002"
)

type PgPool struct {
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func New(config *config.PgConfig, logger *zap.SugaredLogger) (*PgPool, *e.ErrorData) {

	poolConf, err := pgxpool.ParseConfig(config.ToConnString())
	if err != nil {
		return nil, e.NewErrorData(
			DbConnStringParsingError,
			err.Error(),
			fmt.Sprintf("%s.New", currentPackage),
			nil,
			nil,
		)
	}
	poolConf.MinConns = config.MinConns
	poolConf.MaxConns = config.MaxConns
	poolConf.MaxConnIdleTime = config.MaxConnIdleTime
	poolConf.MaxConnLifetime = config.MaxConnLifetime
	poolConf.MaxConnLifetimeJitter = config.MaxConnLifetimeJitter
	poolConf.HealthCheckPeriod = config.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConf)
	if err != nil {
		return nil, e.NewErrorData(
			DbTxPoolCreatingError,
			err.Error(),
			fmt.Sprintf("%s.New", currentPackage),
			nil,
			nil,
		)
	}
	logger.Debugf("pgx connection pool initialized on %s", poolConf.ConnString())
	return &PgPool{
		pool:   pool,
		logger: logger,
	}, nil
}

func (p *PgPool) ShutDownPool() {
	p.pool.Close()
	p.logger.Debug("pgx connection pool shutted down")
}

func (p *PgPool) GetPool() *pgxpool.Pool {
	return p.pool
}
