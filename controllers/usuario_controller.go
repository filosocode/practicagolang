package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/filosocode/practicagolang/auth"
	"github.com/filosocode/practicagolang/data"
	"github.com/filosocode/practicagolang/models"
	"github.com/filosocode/practicagolang/utils"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ----------- CRUD USUARIOS -----------

func GetUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var usuarios []models.Usuario
	if err := data.DB.Preload("Rol").Find(&usuarios).Error; err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}

	var response []models.UsuarioResponse
	for _, u := range usuarios {
		response = append(response, u.ToResponse())
	}

	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Lista de usuarios",
		StatusCode: http.StatusOK,
		Data:       response,
	})
}

func GetUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	var usuario models.Usuario
	if err := data.DB.Preload("Rol").First(&usuario, id).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Usuario encontrado",
		StatusCode: http.StatusOK,
		Data:       usuario.ToResponse(),
	})
}

func NewUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var usuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// ❌ Validar si ya existe el usuario
	var existente models.Usuario
	if err := data.DB.Where("correo = ? OR nombre = ?", usuario.Correo, usuario.Nombre).First(&existente).Error; err == nil {
		http.Error(w, "El usuario ya existe con ese nombre o correo", http.StatusConflict)
		return
	}

	// ✅ Validar rol
	var rol models.Rol
	if err := data.DB.First(&rol, usuario.RolId).Error; err != nil {
		http.Error(w, "El rol asignado no existe", http.StatusBadRequest)
		return
	}

	usuario.Prepare()

	if err := data.DB.Create(&usuario).Error; err != nil {
		http.Error(w, "No se pudo crear el usuario", http.StatusInternalServerError)
		return
	}

	data.DB.Preload("Rol").First(&usuario, usuario.ID)

	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Usuario creado exitosamente",
		StatusCode: http.StatusCreated,
		Data:       usuario.ToResponse(),
	})
}

func UpdateUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	var usuario models.Usuario
	if err := data.DB.First(&usuario, id).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	var input models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// ⚠️ Validar si está intentando usar correo/nombre que ya existen (en otro usuario)
	var duplicado models.Usuario
	if err := data.DB.Where("id <> ? AND (correo = ? OR nombre = ?)", usuario.ID, input.Correo, input.Nombre).First(&duplicado).Error; err == nil {
		http.Error(w, "Ya existe otro usuario con ese nombre o correo", http.StatusConflict)
		return
	}

	usuario.Nombre = input.Nombre
	usuario.Correo = input.Correo
	usuario.RolId = input.RolId

	if input.Password != "" {
		usuario.Password = input.Password
	}
	usuario.Prepare()

	if err := data.DB.Save(&usuario).Error; err != nil {
		http.Error(w, "Error al actualizar", http.StatusInternalServerError)
		return
	}

	data.DB.Preload("Rol").First(&usuario, usuario.ID)

	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Usuario actualizado",
		StatusCode: http.StatusOK,
		Data:       usuario.ToResponse(),
	})
}

func DeleteUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	var usuario models.Usuario
	if err := data.DB.First(&usuario, id).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	if err := data.DB.Delete(&usuario).Error; err != nil {
		http.Error(w, "No se pudo eliminar", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Usuario eliminado",
		StatusCode: http.StatusOK,
		Data:       nil,
	})
}

// ----------- AUTENTICACIÓN -----------

type Credenciales struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var credenciales Credenciales
	if err := json.NewDecoder(r.Body).Decode(&credenciales); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := data.DB.Where("correo = ?", credenciales.Correo).First(&usuario).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Datos de acceso incorrectos", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error al consultar usuario", http.StatusInternalServerError)
		}
		return
	}

	if err := VerificarPassword(usuario.Password, credenciales.Password); err != nil {
		http.Error(w, "Datos de acceso incorrectos", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerarToken(usuario.Correo)
	if err != nil {
		http.Error(w, "Error al generar el token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Autenticación exitosa",
		StatusCode: http.StatusOK,
		Data:       token,
	})
}

func VerificarPassword(passwordHashed string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(password))
}
