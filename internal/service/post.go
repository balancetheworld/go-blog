package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/errs"
	"gorm.io/gorm"
)

func GetPostDetail(ctx context.Context, slugOrID string, viewerID uint, role constant.Role, viewerRoleID uint) (dto.PostDetailResponse, error) {
	post, err := findPost(ctx, slugOrID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.PostDetailResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"post not found",
		)
	}
	if err != nil {
		return dto.PostDetailResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get post failed",
		)
	}
	// 调用权限判断函数，判断当前用户是否拥有管理这篇帖子的权限
	canManage := canManagePost(post, viewerID, role)

	// 条件：帖子是私有 或者 帖子内容为空 ，并且 用户没有管理权限
	if !canViewPost(post, viewerID, role, viewerRoleID) {
		// 返回空数据 + 404 未找到错误
		return dto.PostDetailResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"post not found",
		)
	}
	if resolvedPostVisibility(post) != constant.PostVisibilityPrivate && post.Content != "" {
		if err := repo.IncrementPostViewCount(ctx, post.ID); err != nil {
			return dto.PostDetailResponse{}, errs.NewInternalServer(
				http.StatusInternalServerError,
				"increment post view count failed",
			)
		}
		post.ViewCount++
		post.Heat = post.CalculateHeat()
	}

	result := post.ToDto()

	result.Category = toCategoryResponse(post.Category)

	if canManage {
		draftContent := post.DraftContent
		result.DraftContent = &draftContent
	}

	return result, nil
}

func findPost(
	ctx context.Context,
	slugOrID string,
) (model.Post, error) {
	id, err := strconv.ParseUint(slugOrID, 10, 64)
	if err == nil && id > 0 {
		return repo.GetPostByID(ctx, uint(id))
	}

	return repo.GetPostBySlug(ctx, slugOrID)
}

func canManagePost(
	post model.Post,
	viewerID uint,
	role constant.Role,
) bool {
	if role == constant.RoleAdmin {
		return true
	}

	return role == constant.RoleEditor &&
		viewerID > 0 &&
		post.AuthorID == viewerID
}

func resolvedPostVisibility(post model.Post) constant.PostVisibility {
	if post.IsPrivate {
		return constant.PostVisibilityPrivate
	}

	if post.Visibility == "" {
		return constant.PostVisibilityPublic
	}

	return post.Visibility
}

func canViewPost(
	post model.Post,
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) bool {
	if role == constant.RoleAdmin {
		return true
	}

	if post.Content == "" {
		return canManagePost(post, viewerID, role)
	}

	if viewerID > 0 && post.AuthorID == viewerID {
		return true
	}

	switch resolvedPostVisibility(post) {
	case constant.PostVisibilityPublic:
		return true
	case constant.PostVisibilityRoles:
		if viewerRoleID == 0 {
			return false
		}

		for _, visibleRole := range post.VisibleRoles {
			if visibleRole.ID == viewerRoleID {
				return true
			}
		}
	}

	return false
}

func toCategoryResponse(
	category *model.Category,
) *dto.CategoryResponse {
	if category == nil {
		return nil
	}

	return &dto.CategoryResponse{
		ID:          uint64(category.ID),
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func toLabelResponse(label model.Label) dto.LabelResponse {
	return dto.LabelResponse{
		ID:   uint64(label.ID),
		Name: label.Name,
		Slug: label.Slug,
	}
}

func ListLabels(ctx context.Context) ([]dto.LabelResponse, error) {
	labels, err := repo.ListLabels(ctx)
	if err != nil {
		return nil, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list labels failed",
		)
	}

	result := make([]dto.LabelResponse, 0, len(labels))
	for _, label := range labels {
		result = append(result, toLabelResponse(label))
	}

	return result, nil
}

func toPostListItemResponse(
	post model.Post,
) dto.PostListItemResponse {
	return post.ToDtoWithShortContent()
}

func ListPosts(
	ctx context.Context,
	req dto.PostListRequest,
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) (dto.PostListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	filter := repo.PostListFilter{
		Offset:     (req.Page - 1) * req.PageSize,
		Limit:      req.PageSize,
		Keyword:    req.Keyword,
		Type:       req.Type,
		CategoryID: req.CategoryID,
		LabelID:    req.LabelID,
		AuthorID:   req.AuthorID,
		Status:     "published",
		Sort:       req.Sort,
		PublicOnly: true,
		ViewerID:   viewerID,
		ViewerRoleID: viewerRoleID,
	}
	if role == constant.RoleAdmin {
		filter.PublicOnly = false
	}

	if req.Status == "draft" || req.Status == "all" {
		switch role {
		case constant.RoleAdmin:
			filter.PublicOnly = false
			filter.Status = req.Status
		case constant.RoleEditor:
			if viewerID == 0 {
				return dto.PostListResponse{}, errs.NewForbidden(
					http.StatusForbidden,
					"post list access denied",
				)
			}

			filter.PublicOnly = false
			filter.Status = req.Status
			filter.AuthorID = viewerID
		default:
			return dto.PostListResponse{}, errs.NewForbidden(
				http.StatusForbidden,
				"post list access denied",
			)
		}
	}

	posts, total, err := repo.ListPosts(ctx, filter)
	if err != nil {
		return dto.PostListResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list posts failed",
		)
	}

	items := make([]dto.PostListItemResponse, 0, len(posts))
	for _, post := range posts {
		items = append(items, toPostListItemResponse(post))
	}

	return dto.PostListResponse{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func GetRandomPost(
	ctx context.Context,
	viewerID uint,
	role constant.Role,
	viewerRoleID uint,
) (dto.PostDetailResponse, error) {
	post, err := repo.GetRandomPost(
		ctx,
		role != constant.RoleAdmin,
		viewerID,
		viewerRoleID,
	)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.PostDetailResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"post not found",
		)
	}
	if err != nil {
		return dto.PostDetailResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get random post failed",
		)
	}

	result := post.ToDto()
	result.Category = toCategoryResponse(post.Category)

	return result, nil
}

var invalidPostSlugCharacters = regexp.MustCompile(`[^\p{L}\p{N}]+`)

func normalizePostSlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = invalidPostSlugCharacters.ReplaceAllString(value, "-")
	return strings.Trim(value, "-")
}

func generatePostSlug(
	ctx context.Context,
	title string,
) (string, error) {
	base := normalizePostSlug(title)
	if base == "" {
		base = "post"
	}

	if _, err := strconv.ParseUint(base, 10, 64); err == nil {
		base = "post-" + base
	}

	runes := []rune(base)
	if len(runes) > 220 {
		base = string(runes[:220])
	}

	for sequence := 1; ; sequence++ {
		slug := base
		if sequence > 1 {
			slug = fmt.Sprintf("%s-%d", base, sequence)
		}

		exists, err := repo.CheckPostSlugExists(ctx, slug, 0)
		if err != nil {
			return "", errs.NewInternalServer(
				http.StatusInternalServerError,
				"check post slug failed",
			)
		}

		if !exists {
			return slug, nil
		}
	}
}

func validatePostCategory(
	ctx context.Context,
	categoryID *uint,
) error {
	if categoryID == nil || *categoryID == 0 {
		return nil
	}

	_, err := repo.GetCategoryByID(ctx, *categoryID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.NewBadRequest(
			http.StatusBadRequest,
			"category not found",
		)
	}
	if err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"get category failed",
		)
	}

	return nil
}

func resolvePostLabels(
	ctx context.Context,
	ids []uint,
) ([]model.Label, error) {
	uniqueIDs := make([]uint, 0, len(ids))
	seen := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			return nil, errs.NewBadRequest(
				http.StatusBadRequest,
				"label not found",
			)
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}

	labels, err := repo.GetLabelsByIDs(ctx, uniqueIDs)
	if err != nil {
		return nil, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get labels failed",
		)
	}
	if len(labels) != len(uniqueIDs) {
		return nil, errs.NewBadRequest(
			http.StatusBadRequest,
			"label not found",
		)
	}

	return labels, nil
}

func resolvePostVisibleRoles(
	ctx context.Context,
	visibility constant.PostVisibility,
	ids []uint,
) ([]model.Role, error) {
	if visibility != constant.PostVisibilityRoles {
		return []model.Role{}, nil
	}

	uniqueIDs := make([]uint, 0, len(ids))
	seen := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		if id == 0 {
			return nil, errs.NewBadRequest(
				http.StatusBadRequest,
				"visible role not found",
			)
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		uniqueIDs = append(uniqueIDs, id)
	}

	if len(uniqueIDs) == 0 {
		return nil, errs.NewBadRequest(
			http.StatusBadRequest,
			"visible role is required",
		)
	}

	roles, err := repo.GetEnabledRolesByIDs(ctx, uniqueIDs)
	if err != nil {
		return nil, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get visible roles failed",
		)
	}
	if len(roles) != len(uniqueIDs) {
		return nil, errs.NewBadRequest(
			http.StatusBadRequest,
			"visible role not found",
		)
	}

	return roles, nil
}

func CreatePost(
        ctx context.Context,
        authorID uint,
        role constant.Role,
        req dto.CreatePostRequest,
  ) (dto.PostDetailResponse, error) {
        if authorID == 0 ||
                (role != constant.RoleEditor && role != constant.RoleAdmin) {
                return dto.PostDetailResponse{}, errs.NewForbidden(
                        http.StatusForbidden,
                        "create post access denied",
                )
        }

		title := strings.TrimSpace(req.Title)
        if title == "" {
                return dto.PostDetailResponse{}, errs.NewBadRequest(
                        http.StatusBadRequest,
                        "post title is required",
                )
        }

		if err := validatePostCategory(ctx, req.CategoryID); err != nil {
				return dto.PostDetailResponse{}, err
		}
		categoryID := req.CategoryID
		if categoryID != nil && *categoryID == 0 {
			categoryID = nil
		}

		labels, err := resolvePostLabels(ctx, req.LabelIDs)
		if err != nil {
			return dto.PostDetailResponse{}, err
		}

		visibility := req.Visibility
		if req.IsPrivate {
			visibility = constant.PostVisibilityPrivate
		}
		if visibility == "" {
			visibility = constant.PostVisibilityPublic
		}
		visibleRoles, err := resolvePostVisibleRoles(
			ctx,
			visibility,
			req.VisibleRoleIDs,
		)
		if err != nil {
			return dto.PostDetailResponse{}, err
		}

		slug, err := generatePostSlug(ctx, title)
        if err != nil {
                return dto.PostDetailResponse{}, err
        }

		draftContent := req.DraftContent
		content := ""
		var publishedAt *time.Time
		if req.Publish {
                if strings.TrimSpace(draftContent) == "" {
                        return dto.PostDetailResponse{}, errs.NewBadRequest(
                                http.StatusBadRequest,
                                "published post content is required",
                        )
                }
				content = draftContent
			now := time.Now()
			publishedAt = &now
        }

        postType := strings.TrimSpace(req.Type)
        if postType == "" {
                postType = "article"
        }

        post := model.Post{
                PostBase: model.PostBase{
                        Title:        title,
                        Content:      content,
                        DraftContent: draftContent,
                        Description:  req.Description,
                        Cover:        req.Cover,
                        Type:         postType,
                        Slug:         slug,
						CategoryID:   categoryID,
						IsPrivate:    visibility == constant.PostVisibilityPrivate,
						Visibility:   visibility,
						Top:          req.Top,
						PublishedAt:  publishedAt,
				},
				AuthorID: authorID,
				Labels:   labels,
				VisibleRoles: visibleRoles,
		}

        if err := repo.CreatePost(ctx, &post); err != nil {
                return dto.PostDetailResponse{}, errs.NewInternalServer(
                        http.StatusInternalServerError,
                        "create post failed",
                )
        }

        createdPost, err := repo.GetPostByID(ctx, post.ID)
        if err != nil {
                return dto.PostDetailResponse{}, errs.NewInternalServer(
                        http.StatusInternalServerError,
                        "get created post failed",
                )
        }

        result := createdPost.ToDto()
        result.Category = toCategoryResponse(createdPost.Category)

        draft := createdPost.DraftContent
        result.DraftContent = &draft

        return result, nil
  }

  // UpdatePost 更新帖子业务逻辑
// ctx：上下文，传递请求链路信息、超时、日志、trace等
// id：待修改帖子的数据库主键ID
// viewerID：当前操作人的用户ID，用于权限校验
// role：当前操作人角色（普通用户/管理员/超级管理员）
// req：前端传入的更新帖子请求DTO，所有可修改字段均为指针，区分“前端不传”和“传空值”
// 返回值1：更新完成后的帖子详情对外响应结构体
// 返回值2：业务错误，无错误则为nil
func UpdatePost(
	ctx context.Context,
	id uint,
	viewerID uint,
	role constant.Role,
	req dto.UpdatePostRequest,
) (dto.PostDetailResponse, error) {
	// 1. 根据帖子ID查询数据库完整帖子数据
	post, err := repo.GetPostByID(ctx, id)

	// 分支1：GORM专属错误——这条帖子在数据库中根本不存在
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 返回空响应结构体，自定义404资源不存在错误返回上层handler
		return dto.PostDetailResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"post not found",
		)
	}
	// 分支2：其他查询异常（数据库断开、SQL语法错误、网络故障等）
	if err != nil {
		// 返回空结构体，抛出500服务内部异常，提示查询帖子失败
		return dto.PostDetailResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get post failed",
		)
	}

	// 2. 权限校验：判断当前用户是否具备修改该帖子的权限
	// canManagePost规则：作者本人 / 管理员角色返回true，其余游客、普通用户无权限
	if !canManagePost(post, viewerID, role) {
		// 403 禁止访问，无权限修改帖子
		return dto.PostDetailResponse{}, errs.NewForbidden(
			http.StatusForbidden,
			"update post access denied",
		)
	}

	if req.LabelIDs != nil {
		labels, err := resolvePostLabels(ctx, *req.LabelIDs)
		if err != nil {
			return dto.PostDetailResponse{}, err
		}
		post.Labels = labels
	}

	// 4. 分类ID处理：前端传入了分类ID，校验分类合法性后赋值
	if req.CategoryID != nil {
		// 调用校验函数，判断传入的CategoryID是否存在、是否合法
		if err := validatePostCategory(ctx, req.CategoryID); err != nil {
			// 分类校验失败直接返回错误（如分类不存在、禁用分类）
			return dto.PostDetailResponse{}, err
		}
		// 校验通过，覆盖原帖子的分类ID
		if *req.CategoryID == 0 {
			post.CategoryID = nil
		} else {
			post.CategoryID = req.CategoryID
		}
	}

	// 5. 标题更新逻辑：前端传入标题字段
	if req.Title != nil {
		// 清除标题首尾空格，兼容用户输入多余空格
		title := strings.TrimSpace(*req.Title)
		// 校验：标题清空后为空字符串，不允许空白标题
		if title == "" {
			return dto.PostDetailResponse{}, errs.NewBadRequest(
				http.StatusBadRequest,
				"post title is required",
			)
		}
		// 合法标题赋值给帖子模型
		post.Title = title
	}

	// 7. 草稿内容更新：直接覆盖草稿文本
	if req.DraftContent != nil {
		post.DraftContent = *req.DraftContent
	}

	// 8. 简介更新
	if req.Description != nil {
		post.Description = *req.Description
	}

	// 9. 封面图地址更新
	if req.Cover != nil {
		post.Cover = *req.Cover
	}

	// 10. 帖子类型处理：图文/视频/公告等分类
	if req.Type != nil {
		// 清除首尾空格
		post.Type = strings.TrimSpace(*req.Type)
		// 类型传空时，设置默认类型article文章
		if post.Type == "" {
			post.Type = "article"
		}
	}

	// 11. 是否私有帖子：true仅作者/管理员可见，false公开浏览
	visibility := resolvedPostVisibility(post)
	if req.Visibility != nil {
		visibility = *req.Visibility
	} else if req.IsPrivate != nil {
		if *req.IsPrivate {
			visibility = constant.PostVisibilityPrivate
		} else if visibility == constant.PostVisibilityPrivate {
			visibility = constant.PostVisibilityPublic
		}
	}
	visibleRoleIDs := make([]uint, 0, len(post.VisibleRoles))
	for _, visibleRole := range post.VisibleRoles {
		visibleRoleIDs = append(visibleRoleIDs, visibleRole.ID)
	}
	if req.VisibleRoleIDs != nil {
		visibleRoleIDs = *req.VisibleRoleIDs
	}
	visibleRoles, err := resolvePostVisibleRoles(
		ctx,
		visibility,
		visibleRoleIDs,
	)
	if err != nil {
		return dto.PostDetailResponse{}, err
	}
	post.IsPrivate = visibility == constant.PostVisibilityPrivate
	post.Visibility = visibility
	post.VisibleRoles = visibleRoles

	// 12. 是否置顶帖子
	if req.Top != nil {
		post.Top = *req.Top
	}

	// 13. 发布状态核心逻辑
	if req.Publish != nil {
		// 用户选择发布帖子
		if *req.Publish {
			// 发布必须有正文草稿内容，空白草稿禁止发布
			if strings.TrimSpace(post.DraftContent) == "" {
				return dto.PostDetailResponse{}, errs.NewBadRequest(
					http.StatusBadRequest,
					"published post content is required",
				)
			}
				// 将草稿内容同步到正式展示的Content字段
				post.Content = post.DraftContent
				if post.PublishedAt == nil {
					now := time.Now()
					post.PublishedAt = &now
				}
			} else {
				// 取消发布，清空正式展示内容，仅保留草稿
				post.Content = ""
				post.PublishedAt = nil
		}
	}

	// 14. 调用repo层执行数据库UPDATE更新操作
	if err := repo.UpdatePost(ctx, &post); err != nil {
		// 更新数据库失败，抛出500内部错误
		return dto.PostDetailResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"update post failed",
		)
	}

	// 15. 更新完成后，重新查询一遍最新帖子完整数据（刷新关联分类、最新字段）
	updatedPost, err := repo.GetPostByID(ctx, post.ID)
	if err != nil {
		return dto.PostDetailResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get updated post failed",
		)
	}

	// 16. 数据库模型转换为对外输出DTO，脱敏、过滤敏感字段
	result := updatedPost.ToDto()
	// 拼接分类信息到返回体
	result.Category = toCategoryResponse(updatedPost.Category)

	// 草稿内容赋值给返回结构体指针字段
	draft := updatedPost.DraftContent
	result.DraftContent = &draft

	// 全部逻辑无异常，返回组装完成的帖子详情数据，错误nil
	return result, nil
}

 func DeletePost(
        ctx context.Context,
        id uint,
        viewerID uint,
        role constant.Role,
  ) error {
        post, err := repo.GetPostByID(ctx, id)
        if errors.Is(err, gorm.ErrRecordNotFound) {
                return errs.NewNotFound(
                        http.StatusNotFound,
                        "post not found",
                )
        }
        if err != nil {
                return errs.NewInternalServer(
                        http.StatusInternalServerError,
                        "get post failed",
                )
        }

        if !canManagePost(post, viewerID, role) {
                return errs.NewForbidden(
                        http.StatusForbidden,
                        "delete post access denied",
                )
        }

        rowsAffected, err := repo.DeletePost(ctx, id)
        if err != nil {
                return errs.NewInternalServer(
                        http.StatusInternalServerError,
                        "delete post failed",
                )
        }
        if rowsAffected == 0 {
                return errs.NewNotFound(
                        http.StatusNotFound,
                        "post not found",
                )
        }

        return nil
  }
