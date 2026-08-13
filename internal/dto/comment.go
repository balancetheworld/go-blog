package dto

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
)

type CreateCommentRequest struct {
	PostID     uint                `json:"post_id"`
	TargetType constant.TargetType `json:"target_type" vd:"$ == '' || in($, 'post', 'page', 'comment', 'diary')"`
	TargetID   uint                `json:"target_id"`
	ParentID   *uint               `json:"parent_id"`
	Content    string              `json:"content" vd:"len($) >= 1 && len($) <= 2000"`
}

type CommentListRequest struct {
	PostID     uint                `query:"post_id"`
	TargetType constant.TargetType `query:"target_type" vd:"$ == '' || in($, 'post', 'page', 'diary')"`
	TargetID   uint                `query:"target_id"`
	Page       int                 `query:"page" vd:"$ == 0 || $ >= 1"`
	PageSize   int                 `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
}

type AdminCommentListRequest struct {
	TargetType constant.TargetType `query:"target_type" vd:"$ == '' || in($, 'post', 'page', 'diary')"`
	TargetID   uint                `query:"target_id"`
	AuthorID   uint                `query:"author_id"`
	Keyword    string              `query:"keyword" vd:"len($) <= 100"`
	Page       int                 `query:"page" vd:"$ == 0 || $ >= 1"`
	PageSize   int                 `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
}

type CommentResponse struct {
	ID          uint64              `json:"id"`
	PostID      uint64              `json:"post_id,omitempty"`
	TargetType  constant.TargetType `json:"target_type"`
	TargetID    uint64              `json:"target_id"`
	ParentID    *uint64             `json:"parent_id"`
	RootID      *uint64             `json:"root_id"`
	ReplyToUser *UserDto            `json:"reply_to_user,omitempty"`
	Content     string              `json:"content"`
	Author      UserDto             `json:"author"`
	Depth       uint8               `json:"depth"`
	ReplyCount  uint64              `json:"reply_count"`
	LikeCount   uint64              `json:"like_count"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type CommentListResponse struct {
	Items    []CommentResponse `json:"items"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}
