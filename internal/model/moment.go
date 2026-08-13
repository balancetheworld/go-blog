package model

import "gorm.io/gorm"

type Moment struct {
	gorm.Model
	Content  string `gorm:"type:text;not null"`
	AuthorID uint   `gorm:"not null;index"`
	Author   User   `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
