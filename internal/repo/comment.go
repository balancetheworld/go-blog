package repo

import (
	"context"
	"errors"

	"github.com/zyj/my-blog/internal/model"
	"gorm.io/gorm"
)

func CreateComment(ctx context.Context, comment *model.Comment) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		result := tx.Model(&model.Post{}).
			Where("id = ?", comment.PostID).
			Updates(map[string]any{
				"comment_count": gorm.Expr("comment_count + ?", 1),
				"heat":          gorm.Expr("heat + ?", 5),
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

func GetCommentByID(ctx context.Context, id uint) (model.Comment, error) {
	if db == nil {
		return model.Comment{}, errors.New("database is not initialized")
	}

	var comment model.Comment
	err := db.WithContext(ctx).
		Preload("Author").
		Preload("Post").
		First(&comment, id).
		Error

	return comment, err
}

func ListCommentsByPostID(
	ctx context.Context,
	postID uint,
	offset int,
	limit int,
) ([]model.Comment, int64, error) {
	if db == nil {
		return nil, 0, errors.New("database is not initialized")
	}

	query := db.WithContext(ctx).
		Model(&model.Comment{}).
		Where("post_id = ?", postID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var comments []model.Comment
	err := query.
		Preload("Author").
		Order("created_at DESC, id DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).
		Error

	return comments, total, err
}

func DeleteComment(ctx context.Context, comment model.Comment) (int64, error) {
	if db == nil {
		return 0, errors.New("database is not initialized")
	}

	var rowsAffected int64
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&model.Comment{}, comment.ID)
		if result.Error != nil {
			return result.Error
		}

		rowsAffected = result.RowsAffected
		if rowsAffected == 0 {
			return nil
		}

		postResult := tx.Model(&model.Post{}).
			Where("id = ?", comment.PostID).
			Updates(map[string]any{
				"comment_count": gorm.Expr(
					"CASE WHEN comment_count > 0 THEN comment_count - 1 ELSE 0 END",
				),
				"heat": gorm.Expr(
					"CASE WHEN heat >= 5 THEN heat - 5 ELSE 0 END",
				),
			})
		if postResult.Error != nil {
			return postResult.Error
		}
		if postResult.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})

	return rowsAffected, err
}
