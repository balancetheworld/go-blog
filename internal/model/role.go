package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Code          string `gorm:"size:64;not null;uniqueIndex"`
	Name          string `gorm:"size:64;not null"`
	Description   string `gorm:"size:512;not null;default:''"`
	IsSystem      bool   `gorm:"not null;default:false"`
	IsDefault     bool   `gorm:"not null;default:false"`
	IsRequestable bool   `gorm:"not null;default:false"`
	Enabled       bool   `gorm:"not null;default:true"`
}
