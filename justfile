# Justfile for local development

@default:
  just --list

# Open a psql shell to the running container
psql:
	PGPASSWORD=pass psql -h db -d translatehub

# Run DB migrations (apply schema changes)
migrate-db:
	PGPASSWORD=pass psql -h db -d translatehub < backend/db_schema.sql

build-go:
	cd backend && go build -o ../bin/translatehub  && cd ..

run-go:
	cd backend && go run main.go & 

test-go:
	cd backend && go test ./...

run-app:
	cd frontend && npm start

deploy-local:
	just run-go
	just run-app
