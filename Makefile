postgres:
	docker run -d -p 5432:5432 --name udm-ts-postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret123 postgres
createdb:
	docker exec -it udm-ts-postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it udm-ts-postgres dropdb simple_bank
	
migrateup:
	migrate -path db/migration -database "postgres://root:secret123@localhost:5432/simple_bank?sslmode=disable" -verbose up
	
migratedown:
	migrate -path db/migration -database "postgres://root:secret123@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test