package model

import "gorm.io/gorm"

type Category struct {
        gorm.Model
        Name        string `gorm:"size:64;not null;uniqueIndex"`
        Slug        string `gorm:"size:64;not null;uniqueIndex"`
        Description string `gorm:"size:512;not null;default:''"`
  }
