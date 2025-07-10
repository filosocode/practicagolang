package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/filosocode/practicagolang/data"
	"github.com/filosocode/practicagolang/models"
	"github.com/filosocode/practicagolang/utils"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GET /usuarios
func GetUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var usuarios []models.Usuario
	if err := data.DB.Preload("Rol").Find(&usuarios).Error; err != nil {
		http.Error(w, "Error al obtener usuarios", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Lista de usuarios",
		StatusCode: http.StatusOK,
		Data:       usuarios,
	})
}

// GET /usuarios/{id}
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
		Data:       usuario,
	})
}

// POST /usuarios
func NewUsuario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var usuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}
	usuario.Prepare()
	if err := data.DB.Create(&usuario).Error; err != nil {
		http.Error(w, "No se pudo crear el usuario", http.StatusInternalServerError)
		return
	}

	if err := data.DB.Preload("Rol").First(&usuario, usuario.ID).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(utils.Respuesta{
			Msg:        "Error al cargar Rol",
			StatusCode: http.StatusInternalServerError,
			Data:       err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Usuario creado",
		StatusCode: http.StatusCreated,
		Data:       usuario,
	})
}

// PUT /usuarios/{id}
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

	usuario.Nombre = input.Nombre
	usuario.Correo = input.Correo
	usuario.RolId = input.RolId
	if input.Password != "" {
		usuario.Password = input.Password // Hasheado en BeforeSave
	}
	usuario.Prepare()

	if err := data.DB.Save(&usuario).Error; err != nil {
		http.Error(w, "Error al actualizar", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Usuario actualizado",
		StatusCode: http.StatusOK,
		Data:       usuario,
	})
}

// DELETE /usuarios/{id}
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

type Credenciales struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

type Claims struct {
	Correo string `json:"correo"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var credenciales Credenciales
	if err := json.NewDecoder(r.Body).Decode(&credenciales); err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Correo: usuario.Correo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("API_SECRET")))
	if err != nil {
		http.Error(w, "Error al crear el token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(utils.Respuesta{
		Msg:        "Autenticación exitosa",
		StatusCode: http.StatusOK,
		Data:       tokenString,
	})
}

func VerificarPassword(passwordHashed string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(password))

}
