package service

  import (
        "context"
        "net/http"
        "strings"
		"errors"
		"gorm.io/gorm"

        "github.com/zyj/my-blog/internal/dto"
        "github.com/zyj/my-blog/internal/model"
        "github.com/zyj/my-blog/internal/repo"
        "github.com/zyj/my-blog/pkg/constant"
        "github.com/zyj/my-blog/pkg/errs"
  )

  // ListCategories 查询全部分类列表
// ctx 请求上下文，用于传递链路追踪、超时、日志等信息
// 返回：分类对外展示DTO切片、错误信息
func ListCategories(
	ctx context.Context,
) ([]dto.CategoryResponse, error) {
	// 调用数据层repo查询数据库所有分类数据
	categories, err := repo.ListCategories(ctx)
	// 数据库查询出现异常（连接失败、SQL错误等）
	if err != nil {
		// 返回nil数据，包装500服务内部错误，提示分类列表查询失败
		return nil, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list categories failed",
		)
	}

	// 初始化返回DTO切片，预分配长度，提升append性能
	items := make([]dto.CategoryResponse, 0, len(categories))
	// 遍历数据库原始分类模型，转换为对外输出的DTO结构体
	for _, category := range categories {
		// 转换数据库model为脱敏后的前端响应结构体
		item := toCategoryResponse(&category)
		// 加入结果集合
		items = append(items, *item)
	}

	// 无异常，返回组装完成的分类列表
	return items, nil
}

// CreateCategory 创建新分类业务逻辑
// ctx 请求上下文
// role 当前操作用户角色，用于权限校验
// req 前端提交的创建分类请求参数DTO
// 返回：创建完成的分类详情DTO、错误信息
func CreateCategory(
	ctx context.Context,
	role constant.Role,
	req dto.CreateCategoryRequest,
) (dto.CategoryResponse, error) {
	// 权限校验：判断当前角色是否拥有创建分类权限（仅管理员可创建）
	if !canManageCategory(role) {
		// 403 禁止访问，无操作权限
		return dto.CategoryResponse{}, errs.NewForbidden(
			http.StatusForbidden,
			"create category access denied",
		)
	}

	// 去除分类名称首尾空格，兼容用户输入多余空格
	name := strings.TrimSpace(req.Name)
	// 校验：清空空格后名称为空，参数非法
	if name == "" {
		// 400 参数错误，分类名称不能为空
		return dto.CategoryResponse{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"category name is required",
		)
	}

	// 处理前端传入的自定义短链接slug
	slugSource := strings.TrimSpace(req.Slug)
	// 如果前端没传slug，自动使用分类名称作为生成源
	if slugSource == "" {
		slugSource = name
	}

	// 标准化处理slug：过滤特殊字符、转小写、替换空格为横线，生成友好URL标识
	slug := normalizePostSlug(slugSource)
	// 标准化后slug为空（名称全是特殊符号无法生成合法链接）
	if slug == "" {
		return dto.CategoryResponse{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"category slug is required",
		)
	}

	// 调用repo层校验数据库是否存在同名 / 同slug分类
	// 参数0代表新增场景，不排除自身ID（编辑分类时会传当前分类ID）
	exists, err := repo.CheckCategoryExists(
		ctx,
		name,
		slug,
		0,
	)
	// 校验查询数据库异常
	if err != nil {
		return dto.CategoryResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"check category failed",
		)
	}
	// 数据库已存在重名/重slug分类，冲突无法创建
	if exists {
		// 409 资源冲突错误
		return dto.CategoryResponse{}, errs.NewConflict(
			http.StatusConflict,
			"category name or slug already exists",
		)
	}

	// 组装数据库模型，准备入库
	category := model.Category{
		Name:        name,        // 清洗后的分类名称
		Slug:        slug,        // 标准化后的友好链接标识
		Description: req.Description, // 分类描述，前端不传则为空
	}

	// 调用repo执行数据库新增操作
	if err := repo.CreateCategory(ctx, &category); err != nil {
		// 数据库插入失败，返回500内部错误
		return dto.CategoryResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"create category failed",
		)
	}

	// 将数据库模型转换为前端响应DTO，脱敏后返回
	result := toCategoryResponse(&category)
	return *result, nil
}

  func canManageCategory(role constant.Role) bool {
        return role == constant.RoleEditor ||
                role == constant.RoleAdmin
  }

   func UpdateCategory(
        ctx context.Context,
        id uint,
        role constant.Role,
        req dto.UpdateCategoryRequest,
  ) (dto.CategoryResponse, error) {
        if !canManageCategory(role) {
                return dto.CategoryResponse{}, errs.NewForbidden(
                        http.StatusForbidden,
                        "update category access denied",
                )
        }

        category, err := repo.GetCategoryByID(ctx, id)
        if errors.Is(err, gorm.ErrRecordNotFound) {
                return dto.CategoryResponse{}, errs.NewNotFound(
                        http.StatusNotFound,
                        "category not found",
                )
        }
        if err != nil {
                return dto.CategoryResponse{}, errs.NewInternalServer(
                        http.StatusInternalServerError,
                        "get category failed",
                )
        }

        name := category.Name
        if req.Name != nil {
                name = strings.TrimSpace(*req.Name)
                if name == "" {
                        return dto.CategoryResponse{}, errs.NewBadRequest(
                                http.StatusBadRequest,
                                "category name is required",
                        )
                }
        }

        slug := category.Slug
        if req.Slug != nil {
                slugSource := strings.TrimSpace(*req.Slug)
                if slugSource == "" {
                        slugSource = name
                }

                slug = normalizePostSlug(slugSource)
                if slug == "" {
                        return dto.CategoryResponse{}, errs.NewBadRequest(
                                http.StatusBadRequest,
                                "category slug is required",
                        )
                }
        }

        exists, err := repo.CheckCategoryExists(
                ctx,
                name,
                slug,
                category.ID,
        )
        if err != nil {
                return dto.CategoryResponse{}, errs.NewInternalServer(
                        http.StatusInternalServerError,
                        "check category failed",
                )
        }
        if exists {
                return dto.CategoryResponse{}, errs.NewConflict(
                        http.StatusConflict,
                        "category name or slug already exists",
                )
        }

        category.Name = name
        category.Slug = slug

        if req.Description != nil {
                category.Description = *req.Description
        }

        if err := repo.UpdateCategory(ctx, &category); err != nil {
                return dto.CategoryResponse{}, errs.NewInternalServer(
                        http.StatusInternalServerError,
                        "update category failed",
                )
        }

        result := toCategoryResponse(&category)
        return *result, nil
  }

  func DeleteCategory(
        ctx context.Context,
        id uint,
        role constant.Role,
  ) error {
        if !canManageCategory(role) {
                return errs.NewForbidden(
                        http.StatusForbidden,
                        "delete category access denied",
                )
        }

        rowsAffected, err := repo.DeleteCategory(ctx, id)
        if err != nil {
                return errs.NewInternalServer(
                        http.StatusInternalServerError,
                        "delete category failed",
                )
        }
        if rowsAffected == 0 {
                return errs.NewNotFound(
                        http.StatusNotFound,
                        "category not found",
                )
        }

        return nil
  }