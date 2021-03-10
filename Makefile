postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root schlacke

dropdb:
	docker exec -it postgres12 dropdb schlacke

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/schlacke?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/schlacke?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/schlacke?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/schlacke?sslmode=disable" -verbose down 1

sqlc:
	docker run --rm -v /$(shell pwd)/db:/db -w //db kjconroy/sqlc generate

test:
	go test -v -cover ./...

happy:
	go run main.go -table rims -file data/rims_good.dat

mock:
	mockgen -package mockdb -destination db/mock/store.go db/sqlc Querier

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test server mock
