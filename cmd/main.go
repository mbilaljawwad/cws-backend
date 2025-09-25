package main

import (
	"context"
	"cws-backend/internal/config"
	"cws-backend/internal/database"
	"cws-backend/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

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
		Host:     appCfg.DBHost,
		Port:     appCfg.DBPort,
		User:     appCfg.DBUser,
		Password: appCfg.DBPassword,
		DBName:   appCfg.DBName,
	}

	dbManager := database.NewDBManager(dbCfg).WithRetries(2)
	if err := dbManager.Connect(context.Background()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	serverCtx, cancelCtx := context.WithCancel(context.Background())
	group, _ := errgroup.WithContext(serverCtx)
	// signal for graceful shutdown
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	log.Println("Initializing server")
	srv := server.NewServer(serverCtx, appCfg, dbManager)
	log.Println("Server initialized")

	// Run server with error group
	log.Println("Starting server with error group")
	group.Go(func() error {
		log.Println("Starting server in error group")
		srv.Start()
		return nil
	})

	group.Go(func() error {
		<-signalCh
		defer cancelCtx()
		log.Println("received signal to shutdown")
		return nil
	})
}
