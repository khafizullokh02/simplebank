PSQL_URL=postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres-container -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:14-alpine

createdb:
	docker exec -it postgres-container createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres-container dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down

execdb:
	docker exec -it postgres-container psql -U postgres -d simple_bank

sqlc:
	sqlc generate

cleandb:
	docker exec -it postgres-container psql -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" ${PSQL_URL}

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown execdb sqlc cleandb test server