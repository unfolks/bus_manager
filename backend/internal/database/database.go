package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"bus-manager/internal/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "bus_manager"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
		getEnv("DB_TIMEZONE", "Asia/Jakarta"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Depot{},
		&models.Bus{},
		&models.Route{},
		&models.Trip{},
		&models.Driver{},
		&models.BusUpgrade{},
		&models.Transaction{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Seed initial data
	if err := seedData(db); err != nil {
		log.Printf("Warning: Failed to seed data: %v", err)
	}

	DB = db
	log.Println("Database connected and migrated successfully")
	return db, nil
}

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// Test the connection
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		log.Println("Redis connected successfully")
	}

	RDB = rdb
	return rdb
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func seedData(db *gorm.DB) error {
	// Seed initial routes for Indonesia (Java island)
	routes := []models.Route{
		{
			Name:        "Jakarta - Bandung",
			Origin:      "Jakarta",
			Destination: "Bandung",
			OriginLat:   -6.2088,
			OriginLng:   106.8456,
			DestLat:     -6.9175,
			DestLng:     107.6191,
			Distance:    150,
			Duration:    180, // 3 hours
			Popularity:  80,
			Type:        "intercity",
			MinBusType:  "normal",
			BaseFare:    50000, // IDR
		},
		{
			Name:        "Jakarta - Surabaya",
			Origin:      "Jakarta",
			Destination: "Surabaya",
			OriginLat:   -6.2088,
			OriginLng:   106.8456,
			DestLat:     -7.2575,
			DestLng:     112.7521,
			Distance:    785,
			Duration:    660, // 11 hours
			Popularity:  90,
			Type:        "interprovince",
			MinBusType:  "high_decker",
			BaseFare:    250000, // IDR
		},
		{
			Name:        "Bandung - Yogyakarta",
			Origin:      "Bandung",
			Destination: "Yogyakarta",
			OriginLat:   -6.9175,
			OriginLng:   107.6191,
			DestLat:     -7.7956,
			DestLng:     110.3695,
			Distance:    400,
			Duration:    360, // 6 hours
			Popularity:  70,
			Type:        "intercity",
			MinBusType:  "normal",
			BaseFare:    120000, // IDR
		},
		{
			Name:        "Surabaya - Malang",
			Origin:      "Surabaya",
			Destination: "Malang",
			OriginLat:   -7.2575,
			OriginLng:   112.7521,
			DestLat:     -7.9797,
			DestLng:     112.6304,
			Distance:    90,
			Duration:    120, // 2 hours
			Popularity:  85,
			Type:        "intercity",
			MinBusType:  "normal",
			BaseFare:    35000, // IDR
		},
		{
			Name:        "Yogyakarta - Surakarta",
			Origin:      "Yogyakarta",
			Destination: "Surakarta",
			OriginLat:   -7.7956,
			OriginLng:   110.3695,
			DestLat:     -7.5760,
			DestLng:     110.8295,
			Distance:    60,
			Duration:    90, // 1.5 hours
			Popularity:  75,
			Type:        "intercity",
			MinBusType:  "normal",
			BaseFare:    25000, // IDR
		},
	}

	for _, route := range routes {
		var existingRoute models.Route
		err := db.Where("name = ?", route.Name).First(&existingRoute).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&route).Error; err != nil {
				return fmt.Errorf("failed to create route %s: %w", route.Name, err)
			}
			log.Printf("Created route: %s", route.Name)
		}
	}

	return nil
}
