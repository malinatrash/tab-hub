UP := goose postgres "user=postgres password=postgres dbname=tabhub host=localhost port=5432 sslmode=disable" up -dir ./db/migrations

.PHONY: all
all: up run

.PHONY: up
up:
	$(UP)

.PHONY: run
run:
	go run cmd/tabhub/main.go
