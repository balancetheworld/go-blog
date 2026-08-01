package model

import (
        "time"

        "github.com/zyj/my-blog/pkg/constant"
        "gorm.io/gorm"
  )

  type User struct {
        ID           uint64         `gorm:"primaryKey"`
        Username     string         `gorm:"size:32;not null;uniqueIndex"`
        Email        string         `gorm:"size:254;not null;uniqueIndex"`
        PasswordHash string         `gorm:"size:255;not null" json:"-"`
        Nickname     string         `gorm:"size:64;not null;default:''"`
        Avatar       string         `gorm:"size:512;not null;default:''"`
        Role         constant.Role  `gorm:"size:16;not null;default:user"`
        CreatedAt    time.Time
        UpdatedAt    time.Time
        DeletedAt    gorm.DeletedAt `gorm:"index"`
  }