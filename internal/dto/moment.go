package dto

import "time"

type CreateMomentRequest struct {
	Content string `json:"content" vd:"len($) <= 2000"`
}

type ListMomentsRequest struct {
	Page     int `query:"page" vd:"$ == 0 || $ >= 1"`
	PageSize int `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
}

type MomentResponse struct {
	ID        uint64    `json:"id"`
	Content   string    `json:"content"`
	Author    UserDto   `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListMomentsResponse struct {
	Items    []MomentResponse `json:"items"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}
