package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Usuario struct {
	gorm.Model
	ID       uint64 `gorm:"primary_key;autoIncrement" json:"id"`
	Nombre   string `gorm:"unique;not null" json:"nombre"`
	Correo   string `gorm:"size:100;unique;not null" json:"Email"`
	Password string `gorm:"default:true" json:"password"`
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
