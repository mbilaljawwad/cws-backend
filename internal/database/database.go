package database

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:embed migrations/*sql
var migrationsFS embed.FS

type DBConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	DBName         string
	MigrationTable string
}

type DBManager struct {
	cfg        *DBConfig
	db         *sqlx.DB
	maxRetries int
	migration  *migrate.Migrate
}

func NewDBManager(cfg *DBConfig) *DBManager {
	return &DBManager{
		cfg:        cfg,
		maxRetries: 3,
	}
}

func (dm *DBManager) Connect(ctx context.Context) error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dm.cfg.Host, dm.cfg.Port, dm.cfg.User, dm.cfg.Password, dm.cfg.DBName)
	fmt.Println(connStr)
	var err error
	// retry to connect to database
	for retries := 0; retries < dm.maxRetries; retries++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// initialize database connection
			dm.db, err = sqlx.Connect("postgres", connStr)
			if err != nil {
				log.Printf("failed to connect to database with retries %d : %v", retries, err)
				time.Sleep(time.Duration(retries+1) * time.Second)
				continue
			}

			log.Println("Successfully connected to database")
			return nil
		}
	}
	return fmt.Errorf("failed to connect to database after %d retries", dm.maxRetries)
}

func (dm *DBManager) WithRetries(retries int) *DBManager {
	dm.maxRetries = retries
	return dm
}

func (dm *DBManager) Close() error {
	if dm.db != nil {
		return dm.db.Close()
	}
	return nil
}

func (dm *DBManager) InitMigrator() error {
	if dm.db == nil {
		return fmt.Errorf("database connection is not established")
	}

	// create instance of postgres driver
	postgresCfg := &postgres.Config{
		MigrationsTable: dm.cfg.MigrationTable,
	}
	dbDriver, err := postgres.WithInstance(dm.db.DB, postgresCfg)
	if err != nil {
		return fmt.Errorf("failed to create database driver: %v", err)
	}

	// create source driver for embedded migrations
	srcDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source driver: %v", err)
	}

	dm.migration, err = migrate.NewWithInstance("iofs", srcDriver, dm.cfg.DBName, dbDriver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %v", err)
	}
	return nil
}

func (dm *DBManager) RunMigrationsUp() error {
	if dm.migration == nil {
		return errors.New("migration instance is not initialized")
	}

	if err := dm.migration.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrations to run")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	return nil
}

func (dm *DBManager) GetDatabaseVersion() (uint, bool, error) {
	if dm.migration == nil {
		return 0, false, errors.New("migrator instance is not initialized")
	}

	version, dirty, err := dm.migration.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get database version: %v", err)
	}
	return version, dirty, nil
}
