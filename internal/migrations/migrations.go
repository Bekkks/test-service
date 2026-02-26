package migrations

import (
	"log"

	"subscription-service/internal/models"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	log.Println("Running migrations...")

	// Создание расширения для UUID, если его нет
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		log.Printf("Note: uuid-ossp extension might already exist: %v", err)
	}

	// Автоматическая миграция схемы
	if err := db.AutoMigrate(&models.Subscription{}); err != nil {
		return err
	}

	// Создание индексов, если их нет
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_subscriptions_start_date ON subscriptions(start_date)",
		"CREATE INDEX IF NOT EXISTS idx_subscriptions_end_date ON subscriptions(end_date)",
		"CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions(service_name)",
	}

	for _, idx := range indexes {
		if err := db.Exec(idx).Error; err != nil {
			log.Printf("Index creation note (might already exist): %v", err)
		}
	}

	log.Println("Migrations completed")
	return nil
}
