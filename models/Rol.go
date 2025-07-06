package models

import "gorm.io/gorm"

type Rol struct {
	gorm.Model

	Nombre   string    `gorm:"unique;not null" json:"nombre"`
	Activo   bool      `gorm:"default:true" json:"activo"`
	Usuarios []Usuario `gorm:"default:true" json:"Usuarios"`
}

func (Rol) TableName() string {
	return "roles"

}
