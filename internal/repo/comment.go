package repo

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/pkg/constant"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommentListFilter struct {
	TargetType       constant.TargetType
	TargetID         uint
	AuthorID         uint
	Keyword          string
	ModerationStatus constant.ModerationStatus
	TopLevelOnly     bool
	Offset           int
	Limit            int
	NewestFirst      bool
}

func CreateComment(ctx context.Context, comment *model.Comment) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Author", "ReplyToUser").Create(comment).Error; err != nil {
			return err
		}
		task := model.AITask{
			TaskType:   constant.AITaskCommentModeration,
			TargetType: constant.TargetComment,
			TargetID:   comment.ID,
			Status:     constant.AITaskQueued,
		}
		if err := tx.Create(&task).Error; err != nil {
			return err
		}

		return nil
	})
}

func UpdateCommentModeration(
	ctx context.Context,
	id uint,
	status constant.ModerationStatus,
	reason string,
) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var comment model.Comment
		if err := tx.First(&comment, id).Error; err != nil {
			return err
		}

		if comment.ModerationStatus != constant.ModerationApproved && status == constant.ModerationApproved {
			if err := updateTargetCommentCount(tx, comment.TargetType, comment.TargetID, 1); err != nil {
				return err
			}
			if comment.ParentID != nil {
				if err := updateCommentReplyCount(tx, *comment.ParentID, 1); err != nil {
					return err
				}
			}
		}
		if comment.ModerationStatus == constant.ModerationApproved && status != constant.ModerationApproved {
			if err := updateTargetCommentCount(tx, comment.TargetType, comment.TargetID, -1); err != nil {
				return err
			}
			if comment.ParentID != nil {
				if err := updateCommentReplyCount(tx, *comment.ParentID, -1); err != nil {
					return err
				}
			}
		}

		return tx.Model(&model.Comment{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"moderation_status": status,
				"moderation_reason": reason,
				"moderated_at":      time.Now(),
			}).Error
	})
}

func UpdateCommentModerationResult(
	ctx context.Context,
	id uint,
	status constant.ModerationStatus,
	reason string,
	categories string,
	confidence float64,
) error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var comment model.Comment
		if err := tx.First(&comment, id).Error; err != nil {
			return err
		}

		if comment.ModerationStatus != constant.ModerationApproved && status == constant.ModerationApproved {
			if err := updateTargetCommentCount(tx, comment.TargetType, comment.TargetID, 1); err != nil {
				return err
			}
			if comment.ParentID != nil {
				if err := updateCommentReplyCount(tx, *comment.ParentID, 1); err != nil {
					return err
				}
			}
		}
		if comment.ModerationStatus == constant.ModerationApproved && status != constant.ModerationApproved {
			if err := updateTargetCommentCount(tx, comment.TargetType, comment.TargetID, -1); err != nil {
				return err
			}
			if comment.ParentID != nil {
				if err := updateCommentReplyCount(tx, *comment.ParentID, -1); err != nil {
					return err
				}
			}
		}

		return tx.Model(&model.Comment{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"moderation_status":     status,
				"moderation_reason":     reason,
				"moderation_categories": categories,
				"moderation_confidence": confidence,
				"moderated_at":          time.Now(),
			}).Error
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
	if filter.ModerationStatus != "" {
		query = query.Where("moderation_status = ?", filter.ModerationStatus)
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
	moderationStatus constant.ModerationStatus,
) ([]model.Comment, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	query := db.WithContext(ctx).
		Preload("Author").
		Preload("ReplyToUser").
		Where("parent_id = ?", parentID)
	if moderationStatus != "" {
		query = query.Where("moderation_status = ?", moderationStatus)
	}

	var comments []model.Comment
	err := query.
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

		var approvedCount int64
		if err := query.
			Where("moderation_status = ?", constant.ModerationApproved).
			Count(&approvedCount).Error; err != nil {
			return err
		}

		result := tx.Delete(&model.Comment{}, ids)
		if result.Error != nil {
			return result.Error
		}
		deletedCount = result.RowsAffected

		if comment.ParentID != nil && comment.ModerationStatus == constant.ModerationApproved {
			if err := updateCommentReplyCount(tx, *comment.ParentID, -1); err != nil {
				return err
			}
		}

		return updateTargetCommentCount(
			tx,
			comment.TargetType,
			comment.TargetID,
			-approvedCount,
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

func updateCommentReplyCount(tx *gorm.DB, commentID uint, delta int64) error {
	var expression clause.Expr
	if delta >= 0 {
		expression = gorm.Expr("reply_count + ?", delta)
	} else {
		amount := -delta
		expression = gorm.Expr(
			"CASE WHEN reply_count >= ? THEN reply_count - ? ELSE 0 END",
			amount,
			amount,
		)
	}

	result := tx.Model(&model.Comment{}).
		Where("id = ?", commentID).
		UpdateColumn("reply_count", expression)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
