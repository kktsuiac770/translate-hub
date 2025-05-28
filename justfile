# Justfile for local development

set dotenv-load := true

@default:
  just --list

# Open a psql shell to the running container
psql:
	PGPASSWORD=pass psql -h db -d translatehub

migrate-db:
	migrate -path backend/migrations -database "postgres://vscode:pass@db:5432/translatehub?sslmode=disable" up

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
