package main

import (
	"log"
	"net/http"

	"github.com/filosocode/practicagolang/data"
	"github.com/filosocode/practicagolang/routes"
)

// main es el punto de entrada de la aplicación.
// Inicializa el enrutador y lanza el servidor HTTP en el puerto 8080.
func main() {
	data.ConectarPostgres()
	// Se obtiene el enrutador con todas las rutas configuradas
	rutas := routes.InitRouter()

	// Inicia el servidor en el puerto 8080 y utiliza el enrutador definido
	// Si ocurre un error al iniciar el servidor, se registrará y finalizará la ejecución
	log.Fatal(http.ListenAndServe(":8080", rutas))
}
