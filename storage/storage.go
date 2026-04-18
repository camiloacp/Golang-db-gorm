package storage

import (
	"fmt"
	"log"
	"os"
	"sync"

	"go-db-gorm/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBEngine uint8

const (
	PostgreSQL DBEngine = iota + 1
	MySQL
)

var (
	db   *gorm.DB
	once sync.Once
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func New(engine DBEngine) {
	once.Do(func() {
		var (
			dialect gorm.Dialector
			name    string
		)

		switch engine {
		case PostgreSQL:
			dsn := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
				getEnv("DB_HOST", "localhost"),
				getEnv("DB_USER", "golang_db_user"),
				getEnv("DB_PASSWORD", "golang_db_password"),
				getEnv("DB_NAME", "godb"),
				getEnv("DB_PORT", "7530"),
			)
			dialect = postgres.Open(dsn)
			name = "PostgreSQL"

		case MySQL:
			dsn := fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				getEnv("DB_USER", "root"),
				getEnv("DB_PASSWORD", "root"),
				getEnv("DB_HOST", "127.0.0.1"),
				getEnv("DB_PORT", "3306"),
				getEnv("DB_NAME", "mysql-go"),
			)
			dialect = mysql.Open(dsn)
			name = "MySQL"

		default:
			log.Fatalf("storage.New: unsupported engine %d", engine)
		}

		var err error
		db, err = gorm.Open(dialect, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			log.Fatalf("storage.New: opening %s: %v", name, err)
		}

		fmt.Printf("Database %s connected successfully\n", name)
	})
}

func DB() *gorm.DB {
	return db
}

func Migrate() error {
	return db.AutoMigrate(
		&model.Product{},
		&model.InvoiceHeader{},
		&model.InvoiceItem{},
	)
}
