package main

import (
	"context"
	"cws-backend/internal/config"
	"cws-backend/internal/server"
	"cws-backend/pkg/database"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "cws-backend/docs/swagger" // Import generated docs

	"golang.org/x/sync/errgroup"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:9000
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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

	// Run server with error group
	group.Go(func() error {
		srv.Start()
		return nil
	})

	group.Go(func() error {
		<-signalCh
		defer cancelCtx()
		log.Println("received signal to shutdown")
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Fatalf("error group failed: %v", err)
	}
}
