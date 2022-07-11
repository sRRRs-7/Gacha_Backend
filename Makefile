DB_URL="postgresql://root:secret@localhost:5432/gacha?sslmode=disable"

go:
	go run main.go

postgres:
	docker run --name gachapon -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it gachapon createdb --username=root --owner=root gacha

dropdb:
	docker exec -it gachapon dropdb gacha

migrateinit:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/sRRRs-7/GachaPon/db/sqlc Store


.PHONY: go postgres createdb dropdb migrateinit migrateup migratedown migrateup1 migratedown1 sqlc