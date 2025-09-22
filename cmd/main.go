package main

import (
	"context"
	"cws-backend/internal/config"
	"cws-backend/internal/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func service() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	return router
}

/**
* TODOs:
*  - Implement migration system using golang-migrate
*  - Define proper routes and handlers.
 */

func main() {
	// load configurations
	appCfg := config.Load()

	// connect to database
	dbCfg := &database.DBConfig{
		Host:           appCfg.DBHost,
		Port:           appCfg.DBPort,
		User:           appCfg.DBUser,
		Password:       appCfg.DBPassword,
		DBName:         appCfg.DBName,
		MigrationTable: appCfg.DBMigrationsTable,
	}

	serverCtx, serverCancelCtx := context.WithCancel(context.Background())

	dbManager := database.NewDBManager(dbCfg).WithRetries(2)
	if err := dbManager.Connect(context.Background()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbManager.Close()

	if err := dbManager.InitMigrator(); err != nil {
		log.Fatalf("failed to initialize database migrator: %v", err)
	}

	if err := dbManager.RunMigrationsUp(); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	dbVersion, dirty, err := dbManager.GetDatabaseVersion()
	if err != nil {
		log.Fatalf("failed to get database version: %v", err)
	}
	if dirty {
		log.Println("database is dirty, please run migrations")
	}
	log.Printf("database version: %d", dbVersion)

	// create server context
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", appCfg.APPPort),
		Handler: service(),
	}

	// signal for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-signalCh
		defer serverCancelCtx()

		shutdownCtx, cancelShutdownCtx := context.WithTimeout(serverCtx, time.Second*30)
		defer cancelShutdownCtx()
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out. Forcing exit.")
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("failed to shutdown server: %v", err)
		}
	}()

	log.Println("server started on port 8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server: %v", err)
	}

	// Wait for server context to be stopped.
	<-serverCtx.Done()
}
