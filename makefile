# Makefile

DOWN := go run cmd/goose/down/main.go
UP := go run cmd/goose/up/main.go
RUN := go run cmd/tabhub/main.go
DROP := go run cmd/goose/drop/main.go

.PHONY: all
all: down up run

.PHONY: drop
drop:
	$(DROP)

.PHONY: down
down:
	$(DOWN)

.PHONY: up
up:
	$(UP)

.PHONY: run
run:
	$(RUN)