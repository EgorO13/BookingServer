.PHONY: up down unit-tests e2e-tests seed

up:
	docker-compose up -d

down:
	docker-compose down

unit-tests:
	go test ./internal/... -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

e2e-tests:
	go test ./E2E_test/... -v

seed:
	docker-compose exec app /app/seed