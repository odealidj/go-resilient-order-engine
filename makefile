.PHONY: docker-up docker-down tidy run-order run-inventory run-analytic

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down -v

tidy:
	go mod tidy

run-order:
	go run services/order/cmd/main.go

run-inventory:
	go run services/inventory/cmd/main.go

run-analytic:
	go run services/analytic/cmd/main.go
