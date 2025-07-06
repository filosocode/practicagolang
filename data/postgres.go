package data

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var CONNECTION_STRING = "host=localhost user=postgres password=1234 dbname=api_golang port=5433 sslmode=disable TimeZone=America/Bogota"
var DB *gorm.DB

func ConectarPostgres() {
	var err error
	DB, err = gorm.Open(postgres.Open(CONNECTION_STRING), &gorm.Config{})
	if err != nil {
		log.Fatal("Error de conexión a la base de datos:", err)
	}
	log.Println("Conectado a la BD")
}
