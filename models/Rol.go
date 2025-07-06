package models

import "gorm.io/gorm"

type Rol struct {
	gorm.Model
	Nombre string `gorm:"unique;not null" json:"nombre"`
	Estado bool   `gorm:"default:true" json:"estado"`
}

func (Rol) TableName() string {
	return "roles"

}
