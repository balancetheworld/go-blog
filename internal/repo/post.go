package repo

import (
        "context"
        "errors"
        "strings"

        "github.com/zyj/my-blog/internal/model"
        "gorm.io/gorm"
  )

  // PostListFilter 帖子列表查询筛选条件结构体
// 用于接收前端分页、检索、过滤、排序相关所有查询参数，作为数据库查询条件载体
type PostListFilter struct {
	Offset     int    // 分页偏移量，跳过前Offset条数据，计算公式：(page-1)*limit
	Limit      int    // 分页单页条数，控制单次查询返回多少条帖子
	Keyword   string // 检索关键词，用于模糊匹配帖子标题/内容
	Type       string // 帖子类型筛选，区分不同业务分类的帖子（如图文/视频/公告）
	CategoryID uint // 分类ID，筛选指定栏目/分类下的帖子
	LabelID    uint
	AuthorID   uint   // 作者ID，仅查询该用户发布的帖子
	Status     string // 帖子状态筛选，如草稿/审核中/已发布/已下架
	Sort       string // 排序字段与规则，例如 create_time desc、view_num asc
	PublicOnly bool   // 是否只查询公开帖子；true=仅展示对外公开内容，false=包含私有/仅内部可见帖子
}

 func CreatePost(ctx context.Context, post *model.Post) error {
        if db == nil {
                return errors.New("database is not initialized")
        }

        return db.WithContext(ctx).Create(post).Error
  }

	func UpdatePost(ctx context.Context, post *model.Post) error {
        if db == nil {
                return errors.New("database is not initialized")
        }

		return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.
				Omit("Author", "Category", "Labels").
				Save(post).
				Error; err != nil {
				return err
			}

			return tx.Model(post).Association("Labels").Replace(post.Labels)
		})
	  }

  func DeletePost(ctx context.Context, id uint) (int64, error) {
        if db == nil {
                return 0, errors.New("database is not initialized")
        }

		post := model.Post{}
		post.ID = id
		result := db.WithContext(ctx).
			Select("Labels").
			Delete(&post)
        return result.RowsAffected, result.Error
  }

  func GetPostByID(ctx context.Context, id uint) (model.Post, error) {
        if db == nil {
                return model.Post{}, errors.New("database is not initialized")
        }

        var post model.Post
		err := db.WithContext(ctx).
				Preload("Author").
				Preload("Category").
				Preload("Labels").
				First(&post, id).
                Error

        return post, err
  }

  func GetPostBySlug(ctx context.Context, slug string) (model.Post, error) {
        if db == nil {
                return model.Post{}, errors.New("database is not initialized")
        }

        var post model.Post
		err := db.WithContext(ctx).
				Preload("Author").
				Preload("Category").
				Preload("Labels").
				Where("slug = ?", slug).
                First(&post).
                Error

        return post, err
  }

  func CheckPostSlugExists(
        ctx context.Context,
        slug string,
        excludeID uint,
  ) (bool, error) {
        if db == nil {
                return false, errors.New("database is not initialized")
        }

        query := db.WithContext(ctx).
                Model(&model.Post{}).
                Where("slug = ?", slug)

        if excludeID > 0 {
                query = query.Where("id <> ?", excludeID)
        }

        var count int64
        if err := query.Count(&count).Error; err != nil {
                return false, err
        }

        return count > 0, nil
  }

  func ListPosts(ctx context.Context, filter PostListFilter) ([]model.Post, int64, error) {
	if db == nil {
		return nil, 0, errors.New("database is not initialized")
	}
	query := db.WithContext(ctx).Model(&model.Post{})

	if filter.PublicOnly {
		query = query. 
				Where("is_private = ?", false). 
				Where("content <> ?", "") //只查询 content 字段不为空字符串的数据
	}
// 1. strings.TrimSpace(filter.Keyword)：去除关键词前后空格（用户输入多余空格自动清理）
// 短变量 keyword：接收去空格后的关键词；仅当前if块内有效
if keyword := strings.TrimSpace(filter.Keyword); keyword != "" {
    // 拼接模糊查询通配符 %，前后%代表任意字符，实现包含匹配
    pattern := "%" + keyword + "%"
    // 追加查询条件：标题 或 简介 模糊匹配关键词
    query = query.Where(
        "title LIKE ? OR description LIKE ?",
        pattern, // 第一个? 对应title匹配
        pattern, // 第二个? 对应description匹配
    )
}
if filter.Type != "" {
                query = query.Where("type = ?", filter.Type)
        }

		if filter.CategoryID > 0 {
				query = query.Where("category_id = ?", filter.CategoryID)
		}

		if filter.LabelID > 0 {
			query = query.Where(
				"EXISTS (SELECT 1 FROM post_labels WHERE post_labels.post_id = posts.id AND post_labels.label_id = ?)",
				filter.LabelID,
			)
		}

        if filter.AuthorID > 0 {
                query = query.Where("author_id = ?", filter.AuthorID)
        }

        switch filter.Status {
        case "published":
                query = query.Where("content <> ?", "")
        case "draft":
                query = query.Where("content = ?", "")
        }

        var total int64
        if err := query.Count(&total).Error; err != nil {
                return nil, 0, err
        }

        switch filter.Sort {
        case "oldest":
                query = query.Order("top DESC, created_at ASC")
        case "hot":
                query = query.Order("top DESC, heat DESC, id DESC")
        default:
                query = query.Order("top DESC, created_at DESC")
        }

        var posts []model.Post
		err := query.
				Preload("Author").
				Preload("Category").
				Preload("Labels").
				Offset(filter.Offset).
                Limit(filter.Limit).
                Find(&posts).
                Error
        if err != nil {
                return nil, 0, err
        }

        return posts, total, nil
  }

   func GetRandomPost(
        ctx context.Context,
        publicOnly bool,
  ) (model.Post, error) {
        if db == nil {
                return model.Post{}, errors.New("database is not initialized")
        }

        query := db.WithContext(ctx).Model(&model.Post{})

        if publicOnly {
                query = query.
                        Where("is_private = ?", false).
                        Where("content <> ?", "")
        }

        var post model.Post
		err := query.
				Preload("Author").
				Preload("Category").
				Preload("Labels").
				Order("RANDOM()").
                First(&post).
                Error

        return post, err
  }

  func IncrementPostViewCount(ctx context.Context, id uint) error {
        if db == nil {
                return errors.New("database is not initialized")
        }

		return db.WithContext(ctx).
			Model(&model.Post{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"view_count": gorm.Expr("view_count + ?", 1),
				"heat":       gorm.Expr("heat + ?", 1),
			}).
			Error
  }

  func UpdatePostCommentCount(
        ctx context.Context,
        id uint,
        delta int,
  ) error {
        if db == nil {
                return errors.New("database is not initialized")
        }

		if delta >= 0 {
				return db.WithContext(ctx).
						Model(&model.Post{}).
						Where("id = ?", id).
						Updates(map[string]any{
							"comment_count": gorm.Expr("comment_count + ?", delta),
							"heat":          gorm.Expr("heat + ?", delta*5),
						}).
						Error
        }

        amount := -delta

		return db.WithContext(ctx).
			Model(&model.Post{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"comment_count": gorm.Expr(
					"CASE WHEN comment_count >= ? THEN comment_count - ? ELSE 0 END",
					amount,
					amount,
				),
				"heat": gorm.Expr(
					"CASE WHEN heat >= (CASE WHEN comment_count >= ? THEN ? ELSE comment_count END) * 5 THEN heat - (CASE WHEN comment_count >= ? THEN ? ELSE comment_count END) * 5 ELSE 0 END",
					amount,
					amount,
					amount,
					amount,
				),
			}).
			Error
	  }
