package dto

import (
	"time"

	"github.com/zyj/my-blog/pkg/constant"
)

type RoleOptionResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleApplicationUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type RoleApplicationResponse struct {
	ID            uint                           `json:"id"`
	User          RoleApplicationUserResponse    `json:"user"`
	RequestedRole RoleOptionResponse             `json:"requested_role"`
	Status        constant.RoleApplicationStatus `json:"status"`
	ReviewerID    *uint                          `json:"reviewer_id"`
	ReviewedAt    *time.Time                     `json:"reviewed_at"`
	RejectReason  string                         `json:"reject_reason"`
	CreatedAt     time.Time                      `json:"created_at"`
}

type ListRoleApplicationsRequest struct {
	Page     int                            `query:"page" vd:"$ == 0 || $ >= 1"`
	PageSize int                            `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
	Status   constant.RoleApplicationStatus `query:"status" vd:"$ == '' || in($, 'pending', 'approved', 'rejected')"`
}

type ListRoleApplicationsResponse struct {
	Items    []RoleApplicationResponse `json:"items"`
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"page_size"`
}

type RejectRoleApplicationRequest struct {
	Reason string `json:"reason" vd:"len($) >= 1 && len($) <= 512"`
}
