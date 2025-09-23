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

/**
* TODO:
*  - Implement migration Up/Down commands
*  - Store migration versions in a separate table for backup purposes
*  - Add a command to apply migrations to a specific database
 */

/**
Actions:
- create
- up [N]
- down [N]
- force [V]
- goto [V]
- version
*/

func main() {

	action := flag.String("action", "", "Action to perform")
	name := flag.String("name", "", "Name of the migration")
	// steps := flag.Int("steps", 0, "Number of steps to perform")
	// version := flag.String("version", "", "Version to perform")
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
		fmt.Println("Forcing migration")
	case "goto":
		fmt.Println("Go to migration")
	case "version":
		fmt.Println("Getting version")
	default:
		fmt.Println("Invalid action")
	}
}
