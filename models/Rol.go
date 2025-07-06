package models

import "gorm.io/gorm"

type Rol struct {
	gorm.Model

	ID       uint64    `gorm:"primary_key;autoIncrement" json:"id"`
	Nombre   string    `gorm:"unique;not null" json:"nombre"`
	Estado   bool      `gorm:"default:true" json:"estado"`
	Usuarios []Usuario `gorm:"default:true" json:"Usuarios"`
}

func (Rol) TableName() string {
	return "roles"

}
