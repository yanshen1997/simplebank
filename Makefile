postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123456 -d postgres
createdb:
	docker exec -it postgres14 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres14 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgres://root:123456@localhost:5432/simple_bank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server: 
	go run .
mock:
	mockgen -package mockdb -destination ./db/mock/store.go github.com/yanshen1997/simplebank/db/sqlc Store

.PHONY: createdb dropdb postgres migrateup migratedown sqlc test server mock