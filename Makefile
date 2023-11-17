PSQL_URI=postgres://root:secret@localhost:5432/simple_bank?sslmode=disable
MYSQL_URL=mysql://root:secret@localhost:5432:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres14 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres14 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down

execdb:
	docker exec -it postgres14 psql -U root -d simple_bank

sqlc:
	sqlc generate

cleandb:
	docker exec -it postgres14 psql -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" ${PSQL_URI}

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown execdb sqlc cleandb test