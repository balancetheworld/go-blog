package repo

 import (
        "context"
        "errors"

        "github.com/zyj/my-blog/internal/model"
  )

  func CreateCategory(ctx context.Context, category *model.Category) error {
	if db == nil {
                return errors.New("database is not initialized")
        }
		 return db.WithContext(ctx).
                Create(category).
                Error
  }


  func UpdateCategory(
        ctx context.Context,
        category *model.Category,
  ) error {
        if db == nil {
                return errors.New("database is not initialized")
        }

        return db.WithContext(ctx).
                Save(category). //保存模型数据，自动区分新增 / 更新
                Error
  }

   func DeleteCategory(
        ctx context.Context,
        id uint,
  ) (int64, error) {
        if db == nil {
                return 0, errors.New("database is not initialized")
        }

        result := db.WithContext(ctx).
                Unscoped().
                Delete(&model.Category{}, id)

        return result.RowsAffected, result.Error
  }

   func GetCategoryByID(
        ctx context.Context,
        id uint,
  ) (model.Category, error) {
        if db == nil {
                return model.Category{}, errors.New("database is not initialized")
        }

        var category model.Category
        err := db.WithContext(ctx).
                First(&category, id).
                Error

        return category, err
  }

  func GetCategoryBySlug(
        ctx context.Context,
        slug string,
  ) (model.Category, error) {
        if db == nil {
                return model.Category{}, errors.New("database is not initialized")
        }

        var category model.Category
        err := db.WithContext(ctx).
                Where("slug = ?", slug).
                First(&category).
                Error

        return category, err
  }

  func ListCategories(
        ctx context.Context,
  ) ([]model.Category, error) {
        if db == nil {
                return nil, errors.New("database is not initialized")
        }

        var categories []model.Category
        err := db.WithContext(ctx).
                Order("name ASC").
                Find(&categories).
                Error

        return categories, err
  }

  func CheckCategoryExists(
        ctx context.Context,
        name string,
        slug string,
        excludeID uint,
  ) (bool, error) {
        if db == nil {
                return false, errors.New("database is not initialized")
        }

        query := db.WithContext(ctx).
                Model(&model.Category{}).
                Where("name = ? OR slug = ?", name, slug)

        if excludeID > 0 {
                query = query.Where("id <> ?", excludeID)
        }

        var count int64
        if err := query.Count(&count).Error; err != nil {
                return false, err
        }

        return count > 0, nil
  }
