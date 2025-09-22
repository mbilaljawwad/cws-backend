MIGRATION_DIR=./internal/database/migrations
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=cwsdb
DB_MIGRATIONS_TABLE=migrations

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable&x-migrations-table=$(DB_MIGRATIONS_TABLE)

migrate-create:
	migrate create -ext sql -dir $(MIGRATION_DIR) -tz UTC ${FILE_NAME}

migrate-up:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" down

migrate-status:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" version
	
migrate-force:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" force ${VERSION}

migrate-goto:
	migrate -path $(MIGRATION_DIR) -database "$(DB_URL)" goto ${VERSION}
