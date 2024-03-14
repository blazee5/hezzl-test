migrate-postgres:
	goose -dir ./migrations/postgres postgres "host=localhost user=postgres password=password port=5432 dbname=hezzl sslmode=disable" up

migrate-clickhouse:
	goose -dir ./migrations/clickhouse clickhouse "tcp://127.0.0.1:9000" up