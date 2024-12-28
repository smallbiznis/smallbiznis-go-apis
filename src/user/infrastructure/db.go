package infrastructure

import (
	"fmt"
	"os"
	"time"

	"github.com/smallbiznis/go-lib/pkg/env"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Host     string
	User     string
	Password string
	DbName   string
	Port     string
	SslMode  string
	Timezone string
}

func (d Database) String() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		d.Host, d.User, d.Password, d.DbName, d.Port, d.SslMode, d.Timezone,
	)
}

func NewGorm() (*gorm.DB, error) {
	config := Database{
		Host:     env.Lookup("DB_HOST", "127.0.0.1"),
		User:     env.Lookup("DB_USER", "postgres"),
		Password: env.Lookup("DB_PASSWORD", "35411231"),
		DbName:   env.Lookup("DB_NAME", "postgres"),
		Port:     env.Lookup("DB_PORT", "5432"),
		SslMode:  env.Lookup("DB_SSL_MODE", "disable"),
		Timezone: env.Lookup("DB_TIMEZONE", "Asia/Jakarta"),
	}

	// Initialize the GORM DB connection
	db, err := gorm.Open(postgres.Open(config.String()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Optional: Enable detailed logging
	})
	if err != nil {
		return nil, err
	}

	// Register the OpenTelemetry plugin with GORM
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		zap.L().Fatal("failed use plugin otelgorm", zap.Error(err))
		return nil, err
	}

	if os.Getenv("ENV") != "production" {
		db = db.Debug()
	}

	// Get the underlying SQL DB object to configure the connection pool
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Fatal("failed to get the underlying SQL DB object", zap.Error(err))
		return nil, err
	}

	// Configure the connection pool
	sqlDB.SetMaxOpenConns(100)                 // Set the maximum number of open connections
	sqlDB.SetMaxIdleConns(10)                  // Set the maximum number of idle connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Set the maximum connection lifetime
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // Set the maximum connection idle time

	// Test the connection
	err = sqlDB.Ping()
	if err != nil {
		zap.L().Fatal("failed to ping the database", zap.Error(err))
		return nil, err
	}

	zap.L().Info("Database connection successfully configured with connection pooling.")

	return db, nil
}
