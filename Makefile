createdb:
	docker exec -it postgres createdb --username=udit --owner=udit simple_bank

dropdb:
	docker exec -it postgres dropdb --username=udit simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://udit:root@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://udit:root@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test-all:
	go test -v ./...

coverage-all:
	go test -cover ./...

.PHONY: createdb dropdb migrateup migratedown sqlc test coverage