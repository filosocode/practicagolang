package main

import (
	"log"
	"net/http"

	"github.com/filosocode/practicagolang/data"
	"github.com/filosocode/practicagolang/models"
	"github.com/filosocode/practicagolang/routes"
)

func main() {
	data.ConectarPostgres()

	err := data.DB.AutoMigrate(&models.Rol{})
	if err != nil {
		log.Fatal("Error al migrar modelo Rol:", err)
	}

	rutas := routes.InitRouter()
	log.Fatal(http.ListenAndServe(":8080", rutas))
}
