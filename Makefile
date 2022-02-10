POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/pukbot?sslmode=disable'

dev:
	go run cmd/telebot/main.go fromfile

generate_swagger:
	swag init -g ./internal/server/server.go

migrate_up:
	migrate -database ${POSTGRESQL_URL} -path ./pkg/db/migrations/ up 

migrate_down:
	migrate -database ${POSTGRESQL_URL} -path ./pkg/db/migrations/ down

migrate_force_fix:
	migrate -path ./pkg/db/migrations/ -database ${POSTGRESQL_URL} force ${VERSION}

setup:
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest