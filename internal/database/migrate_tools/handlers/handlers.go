package handlers

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*
var migrationFS embed.FS

const migrationsDir = "./handlers/migrations"

type DBMigrator struct {
	migrator *migrate.Migrate
	db       *sqlx.DB
}

func showError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func (dbm *DBMigrator) InitiateMigrator() error {
	if dbm.db == nil {
		return errors.New("database connection is not established")
	}

	log.Println("Initializing migrator")
	srcDriver, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return err
	}

	dbDriver, err := postgres.WithInstance(dbm.db.DB, &postgres.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		return err
	}

	dbm.migrator, err = migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"postgres",
		dbDriver,
	)
	if err != nil {
		return err
	}
	return nil
}

func NewDBMigrator(db *sqlx.DB) *DBMigrator {
	return &DBMigrator{
		db: db,
	}
}

// Create migration handler
func (dbm *DBMigrator) Create(name string) {
	if name == "" {
		showError(errors.New("-name is required"))
		return
	}

	// create version using UTC timestamp.
	version := time.Now().UTC().Format("20060102150405")

	// create filename using regex to replace all non-alphanumeric characters with an underscore.
	regex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	filename := regex.ReplaceAllString(name, "_")

	// create up and down file paths.
	upfile := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.up.sql", version, filename))
	downfile := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.down.sql", version, filename))

	// create directory if it doesn't exist.
	err := os.MkdirAll(migrationsDir, 0755)
	if err != nil {
		showError(err)
		return
	}
	showError(err)

	// up file template body
	upBody := `BEGIN;
		-- Add Up migration here
	COMMIT;
	`

	// down file template body
	downBody := `BEGIN;
		-- Add Down migration here
	COMMIT;
	`

	// write up file
	if err := os.WriteFile(upfile, []byte(upBody), 0o644); err != nil {
		showError(err)
		return
	}

	// write down file
	if err := os.WriteFile(downfile, []byte(downBody), 0o644); err != nil {
		showError(err)
		return
	}

	fmt.Printf("Migration for %s created successfully with version %s\n", name, version)
	fmt.Println("Up migration file: ", upfile)
	fmt.Println("Down migration file: ", downfile)
}

func (dbm *DBMigrator) Up() {
	log.Println("Up migration is initiated")

	if err := dbm.migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No up migrations to run")
			return
		}
		log.Fatalf("failed to run up migration: %v", err)
		return
	}
	log.Println("Up migration is completed")
}

func (dbm *DBMigrator) Down() {
	log.Println("Initiating down migration")

	if err := dbm.migrator.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No down migrations to run")
			return
		}
		log.Fatalf("failed to run down migration: %v", err)
		return
	}
	log.Println("Down migration is completed")
}
