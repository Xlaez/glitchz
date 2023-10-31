start: 
	docker-compose up
stop:
	docker-compose down
run:
	go run main.go
init-swagger:
	swag init 
.PHONY: start, run, init-swagger, stop