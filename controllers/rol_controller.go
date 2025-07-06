package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/filosocode/practicagolang/utils"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	respuesta := utils.Respuesta{
		Msg:        "Listado de Roles",
		StatusCode: 200,
		Data:       "Listado",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)

}

func NewRol(w http.ResponseWriter, r *http.Request) {
	respuesta := utils.Respuesta{
		Msg:        "Nuevo Rol",
		StatusCode: 200,
		Data:       "Nuevo",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)

}

func GetRol(w http.ResponseWriter, r *http.Request) {
	respuesta := utils.Respuesta{
		Msg:        "Buscar",
		StatusCode: 200,
		Data:       "Buscar uno",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)

}

func DeleteRol(w http.ResponseWriter, r *http.Request) {
	respuesta := utils.Respuesta{
		Msg:        "Borrando",
		StatusCode: 200,
		Data:       "Eliminado",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)

}
