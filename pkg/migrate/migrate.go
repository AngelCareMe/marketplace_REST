package migrate

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

func RunMigrations(dsn, migrationsDir string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logrus.WithError(err).Error("Failed to connect to database for migrations")
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize database driver for migrations")
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsDir,
		"postgres",
		driver,
	)
	if err != nil {
		logrus.WithError(err).Error("Failed to initialize migrator")
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.WithError(err).Error("Failed to apply migrations")
		return err
	}

	logrus.Info("Migrations applied successfully")
	return nil
}
