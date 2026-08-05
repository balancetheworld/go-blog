package middleware

import (
	"testing"

	"github.com/zyj/my-blog/pkg/constant"
)

func TestRoleLevel(t *testing.T) {
	roles := []constant.Role{
		constant.RoleGuest,
		constant.RoleUser,
		constant.RoleEditor,
		constant.RoleAdmin,
	}

	for index, role := range roles {
		if level := roleLevel(role); level != index {
			t.Fatalf("unexpected role level for %s: %d", role, level)
		}
	}

	if roleLevel(constant.Role("unknown")) != -1 {
		t.Fatal("unknown role must be rejected")
	}
}
