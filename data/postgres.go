package data

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConectarPostgres() {

	var error error
	DB, error = gorm.Open(postgres.Open(os.Getenv("CONNECTION_STRING")), &gorm.Config{})
	if error != nil {
		log.Fatal("Error de conexión a la base de datos:", error)
	}
	log.Println("Conectado a la BD")
}
