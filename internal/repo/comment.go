package repo

import (
	"context"
	"errors"
	"strings"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentListFilter struct {
	TargetType   constant.TargetType
	TargetID     uint
	AuthorID     uint
	Keyword      string
	TopLevelOnly bool
	Offset       int
	Limit        int
	NewestFirst  bool
}

func CreateComment(ctx context.Context, comment *model.Comment) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Author", "ReplyToUser").Create(comment).Error; err != nil {
			return err
		}

		if comment.ParentID != nil {
			result := tx.Model(&model.Comment{}).
				Where("id = ?", *comment.ParentID).
				UpdateColumn("reply_count", gorm.Expr("reply_count + ?", 1))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}

		return updateTargetCommentCount(tx, comment.TargetType, comment.TargetID, 1)
	})
}

func GetCommentByID(ctx context.Context, id uint) (model.Comment, error) {
	if db == nil {
		return model.Comment{}, errors.New("database is not initialized")
	}

	var comment model.Comment
	err := db.WithContext(ctx).
		Preload("Author").
		Preload("ReplyToUser").
		First(&comment, id).
		Error

	return comment, err
}

func ListComments(
	ctx context.Context,
	filter CommentListFilter,
) ([]model.Comment, int64, error) {
	if db == nil {
		return nil, 0, errors.New("database is not initialized")
	}

	query := db.WithContext(ctx).Model(&model.Comment{})
	if filter.TargetType != "" {
		query = query.Where("target_type = ?", filter.TargetType)
	}
	if filter.TargetID > 0 {
		query = query.Where("target_id = ?", filter.TargetID)
	}
	if filter.AuthorID > 0 {
		query = query.Where("author_id = ?", filter.AuthorID)
	}
	if filter.TopLevelOnly {
		query = query.Where("parent_id IS NULL")
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
		query = query.Where("content LIKE ?", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var comments []model.Comment
	order := "created_at ASC, id ASC"
	if filter.NewestFirst {
		order = "created_at DESC, id DESC"
	}
	err := query.
		Preload("Author").
		Preload("ReplyToUser").
		Order(order).
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&comments).
		Error

	return comments, total, err
}

func ListCommentReplies(
	ctx context.Context,
	parentID uint,
) ([]model.Comment, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	var comments []model.Comment
	err := db.WithContext(ctx).
		Preload("Author").
		Preload("ReplyToUser").
		Where("parent_id = ?", parentID).
		Order("created_at ASC, id ASC").
		Find(&comments).
		Error

	return comments, err
}

func DeleteComment(ctx context.Context, comment model.Comment) (int64, error) {
	if db == nil {
		return 0, errors.New("database is not initialized")
	}

	var deletedCount int64
	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		query := tx.Model(&model.Comment{})
		switch comment.Depth {
		case 0:
			query = query.Where("id = ? OR root_id = ?", comment.ID, comment.ID)
		case 1:
			query = query.Where("id = ? OR parent_id = ?", comment.ID, comment.ID)
		default:
			query = query.Where("id = ?", comment.ID)
		}

		var ids []uint
		if err := query.Pluck("id", &ids).Error; err != nil {
			return err
		}
		if len(ids) == 0 {
			return nil
		}

		result := tx.Delete(&model.Comment{}, ids)
		if result.Error != nil {
			return result.Error
		}
		deletedCount = result.RowsAffected

		if comment.ParentID != nil {
			if err := tx.Model(&model.Comment{}).
				Where("id = ?", *comment.ParentID).
				UpdateColumn(
					"reply_count",
					gorm.Expr("CASE WHEN reply_count > 0 THEN reply_count - 1 ELSE 0 END"),
				).
				Error; err != nil {
				return err
			}
		}

		return updateTargetCommentCount(
			tx,
			comment.TargetType,
			comment.TargetID,
			-deletedCount,
		)
	})

	return deletedCount, err
}

func updateTargetCommentCount(
	tx *gorm.DB,
	targetType constant.TargetType,
	targetID uint,
	delta int64,
) error {
	if targetType == constant.TargetGuestbook {
		return nil
	}

	var target any
	switch targetType {
	case constant.TargetPost, constant.TargetPage:
		target = &model.Post{}
	case constant.TargetDiary:
		target = &model.Diary{}
	default:
		return errors.New("unsupported comment target type")
	}

	var expression clause.Expr
	if delta >= 0 {
		expression = gorm.Expr("comment_count + ?", delta)
	} else {
		amount := -delta
		expression = gorm.Expr(
			"CASE WHEN comment_count >= ? THEN comment_count - ? ELSE 0 END",
			amount,
			amount,
		)
	}

	result := tx.Model(target).
		Where("id = ?", targetID).
		UpdateColumn("comment_count", expression)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	if targetType != constant.TargetDiary {
		return tx.Model(&model.Post{}).
			Where("id = ?", targetID).
			UpdateColumn(
				"heat",
				gorm.Expr("view_count + like_count * 3 + comment_count * 5"),
			).
			Error
	}

	return nil
}
