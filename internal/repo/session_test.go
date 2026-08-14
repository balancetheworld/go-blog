package repo

import (
	"context"
	"testing"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUpdatePasswordAndRevokeSessions(t *testing.T) {
	database, err := gorm.Open(
		sqlite.Open("file:update_password_session_test?mode=memory&cache=shared"),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(&model.User{}, &model.Session{}); err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = database
	t.Cleanup(func() {
		db = previousDB
	})

	user := model.User{
		Username:     "password-user",
		Email:        "password-user@example.com",
		PasswordHash: "old-password-hash",
		Role:         constant.RoleUser,
	}
	if err := database.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	for _, sessionID := range []string{"session-1", "session-2"} {
		if err := database.Create(&model.Session{
			UserID:    user.ID,
			SessionID: sessionID,
		}).Error; err != nil {
			t.Fatal(err)
		}
	}

	if err := UpdatePasswordAndRevokeSessions(
		context.Background(),
		user.ID,
		"new-password-hash",
	); err != nil {
		t.Fatal(err)
	}

	if err := database.First(&user, user.ID).Error; err != nil {
		t.Fatal(err)
	}
	if user.PasswordHash != "new-password-hash" {
		t.Fatalf("unexpected password hash: %s", user.PasswordHash)
	}
	var sessionCount int64
	if err := database.Model(&model.Session{}).
		Where("user_id = ?", user.ID).
		Count(&sessionCount).
		Error; err != nil {
		t.Fatal(err)
	}
	if sessionCount != 0 {
		t.Fatalf("expected sessions to be revoked, got %d", sessionCount)
	}
}
