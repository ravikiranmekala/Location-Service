package models

import "gorm.io/gorm"

type Location struct {
	gorm.Model
	Id        uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string  `json:"name" gorm:"not null;default:null"`
	Address   string  `json:"address" gorm:"not null;default:null"`
	Latitude  float64 `json:"latitude" gorm:"not null"`
	Longitude float64 `json:"longitude" gorm:"not null"`
	Category  string  `json:"category" gorm:"not null;default:null"`
}
