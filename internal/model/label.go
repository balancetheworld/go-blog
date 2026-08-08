package model

import "gorm.io/gorm"

type Label struct {
	gorm.Model
	Name string `gorm:"size:64;not null;uniqueIndex"`
	Slug string `gorm:"size:64;not null;uniqueIndex"`
}
