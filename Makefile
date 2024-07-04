createdb:
	docker exec -it postgres createdb --username=udit --owner=udit simple_bank

dropdb:
	docker exec -it postgres dropdb --username=udit simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://udit:root@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://udit:root@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://udit:root@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://udit:root@127.0.0.1:5432/simple_bank?sslmode=disable" -verbose down 1


sqlc:
	sqlc generate

test-all:
	go test -v ./...

coverage-all:
	go test -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/uditsaurabh/simple_bank/db/sqlc Store

vet:
	go vet ./...

.PHONY: vet createdb dropdb migrateup migratedown sqlc test coverage server gen-mock migrateup1 migratedown1