package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type DBManager struct {
	cfg        *DBConfig
	db         *sqlx.DB
	maxRetries int
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
			dm.db, err = sqlx.Open("postgres", connStr)
			if err != nil {
				log.Printf("failed to connect to database with retries %d : %v", retries, err)
				time.Sleep(time.Duration(retries+1) * time.Second)
				continue
			}

			// ping database to ensure that the connection is established
			if err = dm.db.PingContext(ctx); err != nil {
				log.Printf("failed to ping database with retries %d : %v", retries, err)
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
