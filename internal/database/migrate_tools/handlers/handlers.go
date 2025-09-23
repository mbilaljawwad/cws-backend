package handlers

import (
	"cws-backend/internal/database"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
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
	DBCfg    *database.DBConfig
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

/*
* Run up migration
 */
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

/*
* Run down migration
 */
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

/*
* Run force migration
 */
func (dbm *DBMigrator) Force(version int) {
	log.Println("Initiating force migration")

	if err := dbm.migrator.Force(version); err != nil {
		log.Fatalf("failed to run force migration: %v", err)
		return
	}
	log.Printf("Force migration for version %d is completed\n", version)
}

/*
* Run goto migration
 */
func (dbm *DBMigrator) Goto(version int) {
	log.Println("Initiating goto migration")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbm.DBCfg.User, dbm.DBCfg.Password, dbm.DBCfg.Host, dbm.DBCfg.Port, dbm.DBCfg.DBName)

	gotoCmd := exec.Command(
		"migrate",
		"-path", "internal/database/migrate_tools/handlers/migrations",
		"-database", connStr,
		"goto", fmt.Sprintf("%d", version))
	gotoCmd.Stdout = os.Stdout
	gotoCmd.Stderr = os.Stderr
	if err := gotoCmd.Run(); err != nil {
		log.Fatalf("failed to run goto migration: %v", err)
		return
	}
	log.Printf("Goto migration for version %d is completed\n", version)
}

/*
* Get current migration version
 */
func (dbm *DBMigrator) Version() {
	log.Println("Getting current migration version")

	version, dirty, err := dbm.migrator.Version()
	if dirty {
		log.Printf("Migration with version %d is dirty\n", version)
	}
	if err != nil {
		log.Fatalf("failed to get current migration version: %v", err)
		return
	}
	log.Printf("Current migration version is %d\n", version)
}
