package utils

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("Error: DATABASE_DSN no se encontró. Asegúrate de que .env esté cargado y la variable exista.")
	}
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
}

func Connect() {
	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Error al conectar la base de datos")
	}

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
