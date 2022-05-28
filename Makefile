postgres:
	sudo docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

startpostgres:
	sudo docker start postgres12

createdb:
	sudo docker exec -it postgres12 createdb --username=root --owner=root gubukid

dropdb:
	sudo docker exec -it postgres12 dropdb gubukid

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/gubukid?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/gubukid?sslmode=disable" -verbose down

sqlc:
	sqlc generate

seed:
	go run seeder/seeder.go

.PHONY: postgres startpostgres createdb dropdb migrateup migratedown sqlc