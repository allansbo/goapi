doc:
	swag fmt -g cmd/main.go
	swag init -g cmd/main.go

up:
	docker compose up -d --build

down:
	docker compose down --volumes --remove-orphans