run:
	go run cmd/server/main.go

migrate-up:
	goose -dir ./migration postgres "user=postgres password=postgres dbname=ticket_app_db sslmode=disable" up

migrate-down:
	goose -dir ./migration postgres "user=postgres password=postgres dbname=ticket_app_db sslmode=disable" down

