package repo

import (
	"context"
	"testing"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestEnsureFixedDiaryFolders(t *testing.T) {
	database, err := gorm.Open(
		sqlite.Open("file:diary_folder_repo_test?mode=memory&cache=shared"),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(
		&model.Role{},
		&model.User{},
		&model.DiaryFolder{},
		&model.Diary{},
	); err != nil {
		t.Fatal(err)
	}

	previousDB := db
	db = database
	t.Cleanup(func() {
		db = previousDB
	})

	user := model.User{
		Username:     "diary-owner",
		Email:        "diary-owner@example.com",
		PasswordHash: "password-hash",
		Role:         constant.RoleAdmin,
	}
	if err := database.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	customFolder := model.DiaryFolder{
		Name:       "自定义",
		Slug:       "custom",
		Visibility: constant.PostVisibilityPublic,
	}
	if err := database.Create(&customFolder).Error; err != nil {
		t.Fatal(err)
	}

	diary := model.Diary{
		Title:      "Diary",
		Slug:       "diary",
		AuthorID:   user.ID,
		FolderID:   &customFolder.ID,
		Visibility: constant.PostVisibilityPublic,
	}
	if err := database.Create(&diary).Error; err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := EnsureFixedDiaryFolders(ctx); err != nil {
		t.Fatal(err)
	}
	if err := EnsureFixedDiaryFolders(ctx); err != nil {
		t.Fatal(err)
	}

	folders, err := ListDiaryFolders(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(folders) != len(fixedDiaryFolders) {
		t.Fatalf("expected %d fixed folders, got %d", len(fixedDiaryFolders), len(folders))
	}
	for index, folder := range folders {
		if folder.Slug != fixedDiaryFolders[index].Slug || folder.Name != fixedDiaryFolders[index].Name {
			t.Fatalf("unexpected folder at index %d: %#v", index, folder)
		}
	}

	var customFolderCount int64
	if err := database.Model(&model.DiaryFolder{}).
		Where("id = ?", customFolder.ID).
		Count(&customFolderCount).
		Error; err != nil {
		t.Fatal(err)
	}
	if customFolderCount != 1 {
		t.Fatalf("expected custom folder to remain, got %d", customFolderCount)
	}

	var updatedDiary model.Diary
	if err := database.First(&updatedDiary, diary.ID).Error; err != nil {
		t.Fatal(err)
	}
	if updatedDiary.FolderID != nil {
		t.Fatalf("expected custom folder association to be cleared, got %d", *updatedDiary.FolderID)
	}
}
