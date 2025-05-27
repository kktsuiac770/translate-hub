# Justfile for local development

@default: 

# Start a PostgreSQL container for local development
start-db:
	docker run --name translatehub-postgres -e POSTGRES_USER=youruser -e POSTGRES_PASSWORD=yourpassword -e POSTGRES_DB=translatehub -p 5432:5432 -d postgres:16

# Stop and remove the PostgreSQL container
stop-db:
	docker stop translatehub-postgres || true
	docker rm translatehub-postgres || true

# View logs for the PostgreSQL container
logs-db:
	docker logs -f translatehub-postgres

# Open a psql shell to the running container
psql:
	docker exec -it translatehub-postgres psql -U youruser -d translatehub

# Run DB migrations (apply schema changes)
migrate-db:
	docker exec -i translatehub-postgres psql -U youruser -d translatehub < db_schema.sql
