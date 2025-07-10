package routes

import (
	"github.com/filosocode/practicagolang/controllers"
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {

	rutas := mux.NewRouter()
	api := rutas.PathPrefix("/api").Subrouter()
	api.HandleFunc("", controllers.GetInitRoute).Methods("GET")

	apiRoles := api.PathPrefix("/roles").Subrouter()
	apiRoles.HandleFunc("", controllers.GetRoles).Methods("GET")
	apiRoles.HandleFunc("/{id}", controllers.GetRol).Methods("GET")
	apiRoles.HandleFunc("/{id}", controllers.UpdateRol).Methods("PUT")
	apiRoles.HandleFunc("", controllers.NewRol).Methods("POST")
	apiRoles.HandleFunc("/{id}", controllers.DeleteRol).Methods("DELETE")

	apiUsuarios := api.PathPrefix("/usuarios").Subrouter()
	apiUsuarios.HandleFunc("", controllers.GetUsuarios).Methods("GET")
	apiUsuarios.HandleFunc("/{id}", controllers.GetUsuario).Methods("GET")
	apiUsuarios.HandleFunc("/{id}", controllers.UpdateUsuario).Methods("PUT")
	apiUsuarios.HandleFunc("", controllers.NewUsuario).Methods("POST")
	apiUsuarios.HandleFunc("/{id}", controllers.DeleteUsuario).Methods("DELETE")

	apiAuth := api.PathPrefix("/auth").Subrouter()

	apiAuth.HandleFunc("/loin", controllers.Login).Methods("POST")

	return rutas
}
