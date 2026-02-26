package main

import (
	"log"
	"subscription-service/internal/config"
	"subscription-service/internal/database"
	"subscription-service/internal/handlers"
	"subscription-service/internal/migrations"
	"subscription-service/internal/router"
)

// @title Subscription Service API
// @version 1.0
// @description REST-сервис для агрегации данных об онлайн подписках пользователей
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация базы данных
	db, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Выполнение миграций
	if err := migrations.Run(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Инициализация обработчиков
	subscriptionHandler := handlers.NewSubscriptionHandler(db)

	// Настройка роутера
	r := router.SetupRouter(subscriptionHandler)

	// Запуск сервера
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
