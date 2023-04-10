MIGRATION_DIR = db_migrations
POSTGRESQL_URL = postgres://mascot-user:password@localhost:5000/mascot?sslmode=disable

setup-migrate: ## Install the migrate tool
	go install github.com/golang-migrate/migrate/v4/cmd/migrate

new-migration:
	migrate create -ext sql -dir $(MIGRATION_DIR) $(name)

migration-up:
	migrate -database ${POSTGRESQL_URL} -path db_migrations up

migration-down:
	migrate -database ${POSTGRESQL_URL} -path db_migrations down $(n)