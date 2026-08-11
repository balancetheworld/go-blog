package service

  import (
        "context"
        "net/http"

        "github.com/zyj/my-blog/internal/dto"
        "github.com/zyj/my-blog/internal/repo"
        "github.com/zyj/my-blog/pkg/errs"
  )

  func ListRequestableRoles(
        ctx context.Context,
  ) ([]dto.RoleOptionResponse, error) {
        roles, err := repo.ListRequestableRoles(ctx)
        if err != nil {
                return nil, errs.NewInternalServer(
                        http.StatusInternalServerError,
                        "list requestable roles failed",
                )
        }

        items := make([]dto.RoleOptionResponse, 0, len(roles))
        for _, role := range roles {
                items = append(items, dto.RoleOptionResponse{
                        ID:          role.ID,
                        Code:        role.Code,
                        Name:        role.Name,
                        Description: role.Description,
                })
        }

        return items, nil
  }

  // ListRoleApplications 获取角色申请列表
// ctx：请求上下文，用于传递超时、取消信号
// req：前端传入的列表查询请求DTO，包含分页、状态筛选条件
// 返回：角色申请列表响应DTO、错误信息
func ListRoleApplications(
	ctx context.Context,
	req dto.ListRoleApplicationsRequest,
) (dto.ListRoleApplicationsResponse, error) {
	// 兼容前端传入page=0的情况，页码最小为1
	if req.Page == 0 {
		req.Page = 1
	}
	// 如果前端没有传PageSize，设置默认每页20条
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 调用repo层查询数据库
	// req.Status：按申请状态筛选(pending/approved/rejected)
	// (req.Page-1)*req.PageSize：计算offset偏移量，实现分页跳过前面的数据
	// req.PageSize：本次查询返回多少条记录
	// applications：数据库查询出来的原始模型数据；total：符合条件的总数据条数
	applications, total, err := repo.ListRoleApplications(
		ctx,
		req.Status,
		(req.Page-1)*req.PageSize,
		req.PageSize,
	)
	// 查询数据库出错，返回空响应体，包装500内部服务错误
	if err != nil {
		return dto.ListRoleApplicationsResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list role applications failed",
		)
	}

	// 初始化返回给前端的DTO切片，预分配容量，避免append频繁扩容，提升性能
	items := make(
		[]dto.RoleApplicationResponse,
		0,
		len(applications),
	)

	// 遍历数据库model，把数据库实体转换成对外输出的DTO结构体
	// 做数据脱敏，只返回前端需要的字段，不返回数据库全部字段
	for _, application := range applications {
		items = append(items, dto.RoleApplicationResponse{
			ID: application.ID, // 申请记录ID
			// 嵌套转换申请人信息DTO
			User: dto.RoleApplicationUserResponse{
				ID:       application.User.ID,       // 用户id
				Username: application.User.Username, // 用户名
				Nickname: application.User.Nickname, // 用户昵称
			},
			// 嵌套转换申请的角色信息DTO
			RequestedRole: dto.RoleOptionResponse{
				ID:          application.RequestedRole.ID,          // 角色ID
				Code:        application.RequestedRole.Code,        // 角色编码
				Name:        application.RequestedRole.Name,        // 角色名称
				Description: application.RequestedRole.Description, // 角色描述
			},
			Status:       application.Status,       // 申请状态 pending/approved/rejected
			ReviewerID:   application.ReviewerID,   // 审核人ID
			ReviewedAt:   application.ReviewedAt,   // 审核时间
			RejectReason: application.RejectReason, // 拒绝理由
			CreatedAt:    application.CreatedAt,    // 申请提交时间
		})
	}

	// 组装分页结果返回给上层handler
	// Items：转换完成的申请数据列表
	// Total：全部符合条件总条数，前端用来渲染分页组件
	// Page、PageSize：回传当前页码和页大小给前端
	return dto.ListRoleApplicationsResponse{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
