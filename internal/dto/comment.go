package dto

import "time"

type CreateCommentRequest struct {
	PostID  uint   `json:"post_id" vd:"$ > 0"`
	Content string `json:"content" vd:"len($) >= 1 && len($) <= 2000"`
}

type CommentListRequest struct {
	PostID   uint `query:"post_id" vd:"$ == 0 || $ > 0"`
	Page     int  `query:"page" vd:"$ == 0 || $ >= 1"`
	PageSize int  `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
}

type CommentResponse struct {
	ID        uint64    `json:"id"`
	PostID    uint64    `json:"post_id"`
	Content   string    `json:"content"`
	Author    UserDto   `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentListResponse struct {
	Items    []CommentResponse `json:"items"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}
