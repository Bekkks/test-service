.PHONY: build run test clean swagger docker-up docker-down

# Установка зависимостей
deps:
	go mod download
	go mod tidy

# Генерация Swagger документации
swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g main.go

# Сборка приложения
build: swagger
	go build -o main .

# Запуск приложения
run: swagger
	go run main.go

# Запуск тестов
test:
	go test -v ./...

# Очистка
clean:
	rm -f main
	rm -rf docs/

# Запуск через Docker Compose
docker-up:
	docker-compose up -d

# Остановка Docker Compose
docker-down:
	docker-compose down

# Пересборка и запуск Docker Compose
docker-rebuild:
	docker-compose up -d --build
