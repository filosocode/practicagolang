package models

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Usuario struct {
	gorm.Model
	ID       uint64 `gorm:"primary_key;autoIncrement" json:"id"`
	Nombre   string `gorm:"unique;not null" json:"nombre"`
	Correo   string `gorm:"size:100;unique;not null" json:"Email"`
	Password string `gorm:"default:true" json:"password"`
	RolId    uint64 `json:"rol_id"`
}

func (Usuario) TableName() string {
	return "usuarios"

}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

}

func VerificarPassword(passwordHashed string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHashed), []byte(password))

}

func (u *Usuario) BeforeSave(tx *gorm.DB) error {
	passwordHashed, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(passwordHashed)
	return nil
}

func (u *Usuario) Prepare() {
	u.ID = 0
	u.Nombre = html.EscapeString(strings.ToUpper(strings.TrimSpace(u.Nombre)))
	u.Correo = html.EscapeString(strings.TrimSpace(u.Correo))
}
