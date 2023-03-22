package migrations

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/Walker088/gorealestate/config"
	e "github.com/Walker088/gorealestate/error"
)

const (
	currentPackage          = "github.com/Walker088/gorealestate/migrate"
	NewSqlDbError           = "MR00001"
	NewDbDriverError        = "MR00002"
	NewMigrateInstanceError = "MR00003"
	MigrateError            = "MR00004"

	migrationSrc = "migrations"
)

type SchemaManager struct {
	rootDir string
	config  *config.PgConfig
	migrate *migrate.Migrate
	logger  *zap.SugaredLogger
}

func New(rootDir string, config *config.PgConfig, logger *zap.SugaredLogger) (*SchemaManager, *e.ErrorData) {
	db, err := sql.Open("pgx", config.ToConnString())
	if err != nil {
		logger.Errorf("[migrate] unable to connect to database: %v\n", err)
		return nil, e.NewErrorData(
			NewSqlDbError,
			err.Error(),
			fmt.Sprintf("%s.New", currentPackage),
			nil,
			nil,
		)
	}
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		logger.Errorf("[migrate] error occured on creating migrate.Migrate: %w", err)
		return nil, e.NewErrorData(
			NewDbDriverError,
			err.Error(),
			fmt.Sprintf("%s.New", currentPackage),
			nil,
			nil,
		)
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file:///%s/%s", rootDir, migrationSrc),
		config.DbName,
		driver,
	)
	if err != nil {
		logger.Errorf("[migrate] error occured on creating schema manager: %w", err)
		return nil, e.NewErrorData(
			NewMigrateInstanceError,
			err.Error(),
			fmt.Sprintf("%s.New", currentPackage),
			nil,
			nil,
		)
	}

	return &SchemaManager{
		rootDir: rootDir,
		config:  config,
		migrate: m,
		logger:  logger,
	}, nil
}

func (sm *SchemaManager) Migrate() *e.ErrorData {
	err := sm.migrate.Up()
	if err == nil {
		v, _, _ := sm.migrate.Version()
		sm.logger.Infof("[migrate] migrated to version %d", v)
		return nil
	}
	if err == migrate.ErrNoChange {
		v, _, _ := sm.migrate.Version()
		sm.logger.Infof("[migrate] there is no schema changes, current version: %d", v)
		return nil
	}
	return e.NewErrorData(
		MigrateError,
		err.Error(),
		fmt.Sprintf("%s.Migrate", currentPackage),
		nil,
		nil,
	)
}

func (sm *SchemaManager) RollBack() {
	sm.migrate.Down()
}

func (sm *SchemaManager) Stop() {
	sm.migrate.Close()
	sm.logger.Debug("Stopped SchemaManager")
}
