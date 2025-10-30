
.PHONY: migrate-up
migrate-up:
	@go run main.go migrate up

.PHONY: migrate-down
migrate-down:
	@go run main.go migrate down

.PHONY: migrate-fresh
migrate-fresh:
	@go run main.go migrate fresh

run : 
	go run main.go rest

docs : 
	swag init -g internal/delivery/http/rest.go --parseDependency true --parseInternal

test:
	go test -v -cover -count=1 -failfast ./... -coverprofile="coverage.out"

migration:
	@name=$(name); \
	migrate create -ext sql -dir database/migration $$name