package dto

import (
	"github.com/zyj/my-blog/pkg/constant"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username" vd:"len($) >= 3 && len($) <= 32"`
	Email    string `json:"email" vd:"email($)"`
	Password string `json:"password" vd:"len($) >= 8 && len($) <= 72"`
	Nickname string `json:"nickname,omitempty" vd:"len($) <= 64"`
}

 type UpdateUserRequest struct {
        Username *string `json:"username,omitempty" vd:"$ == nil || (len($) >= 3 && len($) <= 32)"`
        Email    *string `json:"email,omitempty" vd:"$ == nil || email($)"`
        Nickname *string `json:"nickname,omitempty" vd:"$ == nil || len($) <= 64"`
        Avatar   *string `json:"avatar,omitempty"`
  }

  type UpdateUserRoleRequest struct {
	        Role constant.Role `json:"role" vd:"in($, 'user', 'editor', 'admin')"`
  }

  type ListUsersRequest struct {
        Page     int `query:"page" vd:"$ == 0 || $ >= 1"`
        PageSize int `query:"page_size" vd:"$ == 0 || ($ >= 1 && $ <= 100)"`
  }

  type UserResponse struct {
        ID        uint64        `json:"id"`
        Username  string        `json:"username"`
        Nickname  string        `json:"nickname"`
        Avatar    string        `json:"avatar"`
        Role      constant.Role `json:"role"`
        CreatedAt time.Time     `json:"created_at"`
        UpdatedAt time.Time     `json:"updated_at"`
  }

  type UserPrivateResponse struct {
        UserResponse
        Email string `json:"email"`
  }

  type UserListResponse struct {
        Items    []UserResponse `json:"items"`
        Total    int64          `json:"total"`
        Page     int            `json:"page"`
        PageSize int            `json:"page_size"`
  }
