package service

import (
	"testing"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
)

func TestPostManagementRequiresAdmin(t *testing.T) {
	post := model.Post{AuthorID: 1}

	if !canManagePost(post, 2, constant.RoleAdmin) {
		t.Fatal("admin must be allowed to manage posts")
	}
	if canManagePost(post, 1, constant.RoleEditor) {
		t.Fatal("editor must not be allowed to manage posts")
	}
	if canManagePost(post, 1, constant.RoleUser) {
		t.Fatal("member must not be allowed to manage posts")
	}
}
