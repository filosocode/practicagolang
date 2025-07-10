package main

import (
	"log"
	"net/http"

	"github.com/filosocode/practicagolang/data"
	"github.com/filosocode/practicagolang/models"
	"github.com/filosocode/practicagolang/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error al cargar el archivo .env:", err)
	}

	data.ConectarPostgres()

	data.DB.AutoMigrate(&models.Rol{})
	data.DB.AutoMigrate(&models.Usuario{})

	rutas := routes.InitRouter()
	log.Fatal(http.ListenAndServe(":8080", rutas))
}
