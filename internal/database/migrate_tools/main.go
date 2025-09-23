package main

import (
	"context"
	"cws-backend/internal/config"
	"cws-backend/internal/database"
	"cws-backend/internal/database/migrate_tools/handlers"
	"flag"
	"fmt"
	"log"
)

func main() {
	action := flag.String("action", "", "Action to perform")
	name := flag.String("name", "", "Name of the migration")
	version := flag.Int("version", 0, "Version to perform")
	flag.Parse()

	cfg := config.Load()
	dbCfg := &database.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	}

	dbManager := database.NewDBManager(dbCfg).WithRetries(2)
	if err := dbManager.Connect(context.Background()); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbManager.Close()

	dbMigrator := handlers.NewDBMigrator(dbManager.DB)
	dbMigrator.DBCfg = dbCfg

	err := dbMigrator.InitiateMigrator()
	if err != nil {
		log.Fatalf("failed to initiate migrator: %v", err)
		return
	}

	switch *action {
	case "create":
		dbMigrator.Create(*name)
	case "up":
		dbMigrator.Up()
	case "down":
		dbMigrator.Down()
	case "force":
		dbMigrator.Force(int(*version))
	case "goto":
		dbMigrator.Goto(int(*version))
	case "version":
		dbMigrator.Version()
	default:
		fmt.Println("Invalid action")
	}
}
