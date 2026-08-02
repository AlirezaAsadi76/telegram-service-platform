package migrator

import (
	"database/sql"
	"fmt"
	"telegram-service-platform/repository/postgres"

	_ "github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	config     postgres.DBConfig
	migrations *migrate.FileMigrationSource
}

func New(cfg postgres.DBConfig) Migrator {

	return Migrator{
		config: cfg,

		migrations: &migrate.FileMigrationSource{
			Dir: "repository/migrator/migrations",
		},
	}
}

func (m Migrator) Up() error {

	return m.migrate(migrate.Up)
}

func (m Migrator) Down() error {

	return m.migrate(migrate.Down)
}

func (m Migrator) migrate(
	direction migrate.MigrationDirection,
) error {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		m.config.User,
		m.config.Password,
		m.config.Host,
		m.config.Port,
		m.config.Database,
		m.config.SSLMode,
	)

	db, err := sql.Open(
		"pgx",
		dsn,
	)

	if err != nil {
		return fmt.Errorf(
			"open database: %w",
			err,
		)
	}

	defer db.Close()

	count, err := migrate.Exec(
		db,
		"postgres",
		m.migrations,
		direction,
	)

	if err != nil {

		return fmt.Errorf(
			"execute migration: %w",
			err,
		)
	}

	fmt.Printf(
		"migrated %d files\n",
		count,
	)

	return nil
}
