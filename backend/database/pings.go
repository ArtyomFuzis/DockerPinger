package database

import (
	"gorm.io/gorm"
	"time"
)

type Ping struct {
	gorm.Model
	ServiceId uint
	State     bool
	Date      time.Time `gorm:"sort:desc"`
}

type PingedServices struct {
	gorm.Model
	Address string `gorm:"unique;not null"`
	Pings   []Ping `gorm:"foreignKey:ServiceId"`
}
