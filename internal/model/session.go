package model

import "gorm.io/gorm"

type Session struct {
	// `gorm.Model` 是 GORM 内置的**基础内嵌结构体**，嵌入到你的数据库模型 struct 后，自动给数据表增加 4 个通用字段，不用自己重复定义。
	gorm.Model
	UserID    uint   `gorm:"not null; index"`
	SessionID string `gorm:"size:128;not null;uniqueIndex"`
	UserIP    string `gorm:"size:45;not null;default:''"`
	UserAgent string `gorm:"size:512;not null;default:''"`
}
