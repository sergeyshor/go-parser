createdb:
	docker exec -it postgres15 createdb --username=root --owner=root parser

dropdb:
	docker exec -it postgres15 dropdb parser

migrateup:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/parser?sslmode=disable" -verbose up

migratedown:
	migrate -path migrations -database "postgresql://root:secret@localhost:5432/parser?sslmode=disable" -verbose down

run:
	go run cmd/app/main.go

.PHONY: createdb dropdb migrateup migratedown run