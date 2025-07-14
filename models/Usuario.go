package models

import (
	"html"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Usuario struct {
	gorm.Model
	Nombre   string `gorm:"unique;not null" json:"nombre"`
	Correo   string `gorm:"size:100;unique;not null" json:"correo"`
	Password string `json:"password" gorm:"not null"` // solo se usa en el request
	RolId    uint64 `json:"rolId"`
	Rol      Rol    `json:"rol"`
}

// Estructura segura para las respuestas JSON
type UsuarioResponse struct {
	ID        uint      `json:"id"`
	Nombre    string    `json:"nombre"`
	Correo    string    `json:"correo"`
	RolId     uint64    `json:"rolId"`
	Rol       Rol       `json:"rol"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Asegura el nombre correcto de la tabla
func (Usuario) TableName() string {
	return "usuarios"
}

// Hashea la contraseña antes de guardar
func (u *Usuario) BeforeSave(tx *gorm.DB) error {
	passwordHashed, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(passwordHashed)
	return nil
}

// Limpia y normaliza los datos de entrada
func (u *Usuario) Prepare() {
	u.ID = 0
	u.Nombre = html.EscapeString(strings.ToUpper(strings.TrimSpace(u.Nombre)))
	u.Correo = html.EscapeString(strings.TrimSpace(u.Correo))
}

// Devuelve una versión segura del usuario (sin contraseña)
func (u *Usuario) ToResponse() UsuarioResponse {
	return UsuarioResponse{
		ID:        u.ID,
		Nombre:    u.Nombre,
		Correo:    u.Correo,
		RolId:     u.RolId,
		Rol:       u.Rol,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Función auxiliar para hashear
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
