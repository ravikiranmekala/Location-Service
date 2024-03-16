package database

import (
	"fmt"
	"log"
	"os"

	"github.com/ravikiranmekala/MapUp-Backend-assessment/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		// "host.docker.internal",
		// "ravikiran",
		// "ravikiran",
		// "ravikiran",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database! \n", err)
		os.Exit(2)
	}
	log.Println("Connected to database!")
	db.Logger = db.Logger.LogMode(logger.Info)

	log.Println("Migrating database...")
	db.AutoMigrate(&models.Location{})

	DB = Dbinstance{
		Db: db,
	}
}
