package model

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	PostID        *uint               `gorm:"index"`
	TargetType    constant.TargetType `gorm:"size:16;not null;default:post;index:idx_comment_target,priority:1"`
	TargetID      uint                `gorm:"not null;default:0;index:idx_comment_target,priority:2"`
	ParentID      *uint               `gorm:"index"`
	RootID        *uint               `gorm:"index"`
	ReplyToUserID *uint               `gorm:"index"`
	ReplyToUser   *User               `gorm:"foreignKey:ReplyToUserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	AuthorID      uint                `gorm:"not null;index"`
	Author        User                `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Content       string              `gorm:"type:text;not null"`
	Depth         uint8               `gorm:"not null;default:0"`
	ReplyCount    uint64              `gorm:"not null;default:0"`
	LikeCount     uint64              `gorm:"not null;default:0"`
	ModerationStatus     constant.ModerationStatus `gorm:"size:20;not null;default:approved;index"`
	ModerationReason     string                    `gorm:"type:text"`
	ModerationCategories string                    `gorm:"type:text"`
	ModerationConfidence float64
	ModeratedAt          *time.Time
}

func (c *Comment) BeforeCreate(_ *gorm.DB) error {
	if c.TargetType == "" && c.PostID != nil && *c.PostID > 0 {
		c.TargetType = constant.TargetPost
	}
	if c.TargetID == 0 && c.PostID != nil && *c.PostID > 0 {
		c.TargetID = *c.PostID
	}
	if c.ModerationStatus == "" {
		c.ModerationStatus = constant.ModerationPending
	}

	return nil
}
