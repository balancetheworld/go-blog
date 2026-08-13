package model

import (
	"github.com/zyj/my-blog/internal/dto"
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
}

func (c *Comment) BeforeCreate(_ *gorm.DB) error {
	if c.TargetType == "" && c.PostID != nil && *c.PostID > 0 {
		c.TargetType = constant.TargetPost
	}
	if c.TargetID == 0 && c.PostID != nil && *c.PostID > 0 {
		c.TargetID = *c.PostID
	}

	return nil
}

func (c Comment) ToDto() dto.CommentResponse {
	var parentID *uint64
	if c.ParentID != nil {
		value := uint64(*c.ParentID)
		parentID = &value
	}

	var rootID *uint64
	if c.RootID != nil {
		value := uint64(*c.RootID)
		rootID = &value
	}

	var replyToUser *dto.UserDto
	if c.ReplyToUser != nil {
		value := c.ReplyToUser.ToDto()
		replyToUser = &value
	}

	postID := uint64(0)
	if c.PostID != nil {
		postID = uint64(*c.PostID)
	}
	if c.TargetType == constant.TargetPost || c.TargetType == constant.TargetPage {
		postID = uint64(c.TargetID)
	}

	return dto.CommentResponse{
		ID:          uint64(c.ID),
		PostID:      postID,
		TargetType:  c.TargetType,
		TargetID:    uint64(c.TargetID),
		ParentID:    parentID,
		RootID:      rootID,
		ReplyToUser: replyToUser,
		Content:     c.Content,
		Author:      c.Author.ToDto(),
		Depth:       c.Depth,
		ReplyCount:  c.ReplyCount,
		LikeCount:   c.LikeCount,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
