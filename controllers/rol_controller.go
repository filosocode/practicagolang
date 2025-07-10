package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/filosocode/practicagolang/data"
	"github.com/filosocode/practicagolang/models"
	"github.com/filosocode/practicagolang/utils"
	"github.com/gorilla/mux"
)

func GetRoles(w http.ResponseWriter, r *http.Request) {
	var roles []models.Rol
	data.DB.Find(&roles)
	respuesta := utils.Respuesta{
		Msg:        "Listado de Roles",
		StatusCode: 200,
		Data:       roles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)

}

func NewRol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rol models.Rol
	if err := json.NewDecoder(r.Body).Decode(&rol); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(utils.Respuesta{
			Msg:        "Error en los datos enviados",
			StatusCode: http.StatusBadRequest,
			Data:       err.Error(),
		})
		return
	}

	if err := data.DB.Create(&rol).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.Respuesta{
			Msg:        "Error al intentar crear el rol",
			StatusCode: http.StatusInternalServerError,
			Data:       err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Rol creado con éxito",
		StatusCode: http.StatusCreated,
		Data:       rol,
	})
}

func UpdateRol(w http.ResponseWriter, r *http.Request) {
	var rol models.Rol
	w.Header().Set("Content -Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	if err := json.NewDecoder(r.Body).Decode(&rol); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(utils.Respuesta{
			Msg:        "Error al decodificar la solicitud",
			StatusCode: http.StatusBadRequest,
			Data:       err.Error(),
		})
		return
	}

	var rolExistente models.Rol
	if err := data.DB.First(&rolExistente, id).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(utils.Respuesta{
			Msg:        "Rol no encontrado",
			StatusCode: http.StatusNotFound,
			Data:       err.Error(),
		})
		return
	}

	rolExistente.Nombre = rol.Nombre
	rolExistente.Activo = rol.Activo
	if err := data.DB.Save(&rolExistente).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.Respuesta{
			Msg:        "Error al actualizar el rol",
			StatusCode: http.StatusInternalServerError,
			Data:       err.Error(),
		})
		return
	}

	respuesta := utils.Respuesta{
		Msg:        "Rol actualizado con exito",
		StatusCode: http.StatusOK,
		Data:       rolExistente,
	}

	json.NewEncoder(w).Encode(&respuesta)

}

func GetRol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var rol models.Rol
	data.DB.First(&rol, params["id"])
	if rol.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		respuesta := utils.Respuesta{
			Msg:        "Rol no encontrado",
			StatusCode: http.StatusNotFound,
			Data:       nil,
		}
		json.NewEncoder(w).Encode(respuesta)
		return

	}
	data.DB.Model(&rol).Association("Usuarios").Find(&rol.Usuarios)
	respuesta := utils.Respuesta{
		Msg:        "Rol encontrado",
		StatusCode: http.StatusOK,
		Data:       rol,
	}
	json.NewEncoder(w).Encode(respuesta)

}

func DeleteRol(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var rol models.Rol
	data.DB.First(&rol, params["id"])
	if rol.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		respuesta := utils.Respuesta{
			Msg:        "Rol no encontrado",
			StatusCode: http.StatusNotFound,
			Data:       nil,
		}
		json.NewEncoder(w).Encode(respuesta)
		return

	}
	data.DB.Delete(&rol)
	//data.DB.Unscoped().Delete(&rol)
	respuesta := utils.Respuesta{
		Msg:        "Rol Eliminado",
		StatusCode: http.StatusAccepted,
		Data:       rol,
	}
	json.NewEncoder(w).Encode(respuesta)
}
