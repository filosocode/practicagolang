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

	data.DB.AutoMigrate(&models.Rol{})
	data.DB.AutoMigrate(&models.Usuario{})

	rutas := routes.InitRouter()
	log.Fatal(http.ListenAndServe(":8080", rutas))
}
