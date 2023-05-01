start: 
	docker-compose up
stop:
	docker-compose down
run:
	go run cmd/main.go
init-swagger:
	swag init -g cmd/main.go
.PHONY: start, run, init-swagger, stop