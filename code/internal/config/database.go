package config

import (
	"code/internal/models"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := "postgres://postgres:postgres@localhost:5432/golang_db"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	db.AutoMigrate(&models.User{})
	slog.Info("connected to database")
	return db
}
