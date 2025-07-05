package routes

import (
	"github.com/filosocode/practicagolang/controllers"
	"github.com/gorilla/mux"
)

// InitRouter configura y retorna el enrutador principal de la aplicación.
// Aquí se definen las rutas base de la API utilizando el paquete Gorilla Mux.
func InitRouter() *mux.Router {
	// Se crea un nuevo enrutador raíz
	rutas := mux.NewRouter()

	// Se define un subenrutador bajo el prefijo /api para agrupar endpoints de la API
	api := rutas.PathPrefix("/api").Subrouter()

	// Ruta GET /api o /api/ que responde con un mensaje de estado
	api.HandleFunc("", controllers.GetInitRoute).Methods("GET")
	api.HandleFunc("/", controllers.GetInitRoute).Methods("GET")

	// Se retorna el enrutador principal con las rutas configuradas
	return rutas
}
