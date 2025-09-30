# Makefile for cws-backend
migrate_tools = pkg/migrate_tools/main.go
swag_tools = github.com/swaggo/swag/cmd/swag@latest

# .PHONY: help swagger-gen swagger-gen-docker swagger-watch build dev clean

# create migration
migrate-create:
	go run $(migrate_tools) -action create -name $(name)

# run migrations
migrate-up:
	go run $(migrate_tools) -action up

# run migrations down
migrate-down:
	go run $(migrate_tools) -action down

# run migrations force
migrate-force:
	go run $(migrate_tools) -action force -version $(version)

# run migrations goto
migrate-goto:
	go run $(migrate_tools) -action goto -version $(version)

# run migrations version
migrate-version:
	go run $(migrate_tools) -action version

gen-doc:
	go run $(swag_tools) init -g ./cmd/main.go --parseDependency --parseInternal -o ./docs/swagger