package service

import (
	"context"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/zyj/my-blog/internal/dto"
	"github.com/zyj/my-blog/internal/model"
	"github.com/zyj/my-blog/internal/repo"
	"github.com/zyj/my-blog/pkg/constant"
	"github.com/zyj/my-blog/pkg/errs"
	"github.com/zyj/my-blog/pkg/utils"
)

func ListUsers(ctx context.Context, req dto.ListUsersRequest) (dto.UserListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	users, total, err := repo.ListUsers(ctx, (req.Page-1)*req.PageSize, req.PageSize)
	if err != nil {
		return dto.UserListResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"list users failed",
		)
	}

	items := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, toUserResponse(user))
	}

	return dto.UserListResponse{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func GetUser(ctx context.Context, id uint64) (dto.UserResponse, error) {
	user, err := repo.GetUserByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"user not found",
		)
	}
	if err != nil {
		return dto.UserResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get user failed",
		)
	}

	return toUserResponse(user), nil
}

func CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.UserPrivateResponse, error) {
	// 1. 参数校验已经在 controller 层完成，这里不需要重复校验
	// 2. 业务逻辑处理：调用 repo 层，保存用户数据到数据库
	// 3. 返回结果给 controller 层，controller 层再返回给前端
	exists, err := repo.UserExists(ctx, req.Username, req.Email)
	if err != nil {
		//空的**响应 DTO 结构体**。规范：哪怕出错，也要把出参 DTO 占位返回，避免上层调用方结构体解析报错。
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"query user exists failed",
		)
	}
	if exists {
		return dto.UserPrivateResponse{}, errs.NewConflict(
			http.StatusConflict,
			"user already exists",
		)
	}
	//使用 bcrypt 算法给前端传来的明文密码生成加密哈希，自动加盐，安全保存到数据库，防止密码泄露。
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"hash password failed",
		)
	}

	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Nickname:     req.Nickname,
		Role:         constant.RoleUser,
	}

	if err := repo.CreateUser(ctx, &user); err != nil {
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"create user failed",
		)
	}
	return toUserPrivateResponse(user), nil
}

func UpdateUser(ctx context.Context, id uint64, req dto.UpdateUserRequest) (dto.UserPrivateResponse, error) {
	user, err := repo.GetUserByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.UserPrivateResponse{}, errs.NewNotFound(
			http.StatusNotFound,
			"user not found",
		)
	}
	if err != nil {
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get user failed",
		)
	}

	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Nickname != nil {
		user.Nickname = *req.Nickname
	}
	if req.Avatar != nil {
		user.AvatarURL = *req.Avatar
	}

	exists, err := repo.UserExistsExcept(ctx, id, user.Username, user.Email)
	if err != nil {
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"query user exists failed",
		)
	}
	if exists {
		return dto.UserPrivateResponse{}, errs.NewConflict(
			http.StatusConflict,
			"username or email already exists",
		)
	}

	if err := repo.UpdateUser(ctx, &user); err != nil {
		return dto.UserPrivateResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"update user failed",
		)
	}

	return toUserPrivateResponse(user), nil
}

func DeleteUser(ctx context.Context, id uint64) error {
	rowsAffected, err := repo.DeleteUser(ctx, id)
	if err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"delete user failed",
		)
	}
	if rowsAffected == 0 {
		return errs.NewNotFound(
			http.StatusNotFound,
			"user not found",
		)
	}

	return nil
}

// 转换函数确保数据库模型中的密码不会进入响应
// 专门写转换函数，隔离数据库 model 和前端返回 DTO，彻底阻止密码、敏感字段泄露给前端
// model.User 是数据库模型   dto.UserResp 是前端响应 dto
func toUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
			ID: uint64(user.ID),
		Username:  user.Username,
		Nickname:  user.Nickname,
			Avatar:    user.AvatarURL,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toUserPrivateResponse(user model.User) dto.UserPrivateResponse {
	return dto.UserPrivateResponse{
                UserResponse: toUserResponse(user),
                Email:        user.Email,
        }
}

// userLogin 用户登录业务处理函数（service层核心登录逻辑）
// ctx: 请求上下文，用于数据库查询超时控制、链路透传
// req: 前端传入的登录请求DTO，包含账号、密码、记住我、客户端IP、浏览器标识等信息
// 返回：用户脱敏信息+双令牌响应DTO / 自定义业务错误（统一封装http状态码与提示文案）
func UserLogin(ctx context.Context, req *dto.UserLoginReq) (dto.UserAuthResponse, error) {
	// 1. 参数基础校验：前端请求体不能为空
	if req == nil {
		// 返回400错误，提示缺少登录参数
		return dto.UserAuthResponse{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"login request is required",
		)
	}

	// 2. 根据输入的账号（用户名/邮箱）查询数据库用户
	user, err := repo.GetUserByUsernameOrEmail(
		ctx,
		req.Account,
	)
	// 判断错误类型：GORM未查询到用户记录
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 返回401未授权，统一提示账号或密码错误（安全，不区分是账号不存在还是密码错）
		return dto.UserAuthResponse{}, errs.NewUnauthorized(
			http.StatusUnauthorized,
			"invalid account or password",
		)
	}
	// 数据库查询出现未知异常（连接失败、SQL错误等）
	if err != nil {
		// 返回500服务器内部错误
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"get user failed",
		)
	}

	// 3. 密码校验：bcrypt比对加密哈希密码与前端传入明文密码
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), // 数据库存储的加密后的密码哈希
		[]byte(req.Password),      // 前端提交的原始明文密码
	)
	// 比对失败，密码不正确，返回401
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewUnauthorized(
			http.StatusUnauthorized,
			"invalid account or password",
		)
	}

	// 4. 生成全局唯一SessionID，作为本次登录会话唯一标识
	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"generate session failed",
		)
	}

	// 5. 根据用户ID、角色、会话ID、是否记住我，生成 Access + Refresh 双令牌
	tokens, err := utils.GenerateTokenPair(
		user.ID,
		user.Role,
		sessionID,
		req.Remember,
	)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"generate token failed",
		)
	}

	// 6. 组装数据库会话记录model，存入登录设备信息
	session := model.Session{
		UserID:    user.ID,        // 关联登录用户ID
		SessionID: sessionID,      // 本次会话唯一ID，与JWT内Claims对应
		UserIP:    req.UserIP,     // 客户端登录IP，用于安全审计
		UserAgent: req.UserAgent,  // 客户端浏览器/设备标识
	}

	// 7. 将会话记录插入sessions数据表，留存登录状态，用于后续撤销、校验会话
	if err := repo.CreateSession(ctx, &session); err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"create session failed",
		)
	}

	// 8. 组装登录成功返回给前端的响应DTO
	// toUserPrivateResponse：对数据库用户model脱敏，隐藏密码、敏感字段，只返回展示信息
	return dto.UserAuthResponse{
		User:          toUserPrivateResponse(user),
		AccessToken:   tokens.AccessToken,  // 业务访问短时效令牌
		RefreshToken:  tokens.RefreshToken, // 刷新长时效令牌
		AccessMaxAge:  tokens.AccessMaxAge, // Access过期秒数，前端用于cookie过期控制
		RefreshMaxAge: tokens.RefreshMaxAge,// Refresh过期秒数
	}, nil
}


// UserRegister 用户注册业务逻辑（service层注册接口核心函数）
// ctx：请求上下文，透传给repo做数据库超时控制
// req：前端注册请求DTO，包含用户名、邮箱、密码、昵称、记住我、客户端IP、设备标识
// 返回：登录同款用户+双Token响应结构体 / 统一封装的业务错误
func UserRegister(ctx context.Context, req *dto.UserRegisterReq) (dto.UserAuthResponse, error) {
	// 1. 校验入参：前端请求体不能为空
	if req == nil {
		// 返回400错误：请求参数缺失
		return dto.UserAuthResponse{}, errs.NewBadRequest(
			http.StatusBadRequest,
			"register request is required",
		)
	}

	// 2. 读取环境配置，判断系统是否开启注册功能
	// EnvKeyEnableRegister：环境变量控制注册开关，未配置默认true允许注册
	if !utils.GetAsBool(
		constant.EnvKeyEnableRegister,
		true,
	) {
		// 返回403禁止访问：系统关闭注册
		return dto.UserAuthResponse{}, errs.NewForbidden(
			http.StatusForbidden,
			"registration is disabled",
		)
	}

	// 3. 校验用户名是否已被占用
	usernameExists, err := repo.CheckUsernameExists(
		ctx,
		req.Username,
	)
	if err != nil {
		// 查询数据库异常，返回500服务器错误
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"check username failed",
		)
	}
	if usernameExists {
		// 409冲突：用户名重复
		return dto.UserAuthResponse{}, errs.NewConflict(
			http.StatusConflict,
			"username already exists",
		)
	}

	// 4. 校验邮箱是否已注册
	emailExists, err := repo.CheckEmailExists(
		ctx,
		req.Email,
	)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"check email failed",
		)
	}
	if emailExists {
		// 409冲突：邮箱已存在
		return dto.UserAuthResponse{}, errs.NewConflict(
			http.StatusConflict,
			"email already exists",
		)
	}

	if utils.GetAsBool(
		constant.EnvKeyEnableEmailVerify,
		true,
	) {
		if !utils.VerifyEmailCode(
			req.Email,
			req.Code,
		) {
			return dto.UserAuthResponse{}, errs.NewUnauthorized(
				http.StatusUnauthorized,
				"invalid or expired email verification code",
			)
		}
	}

	// 5. 查询数据库总用户数量，用于判断是否是本站第一个注册用户
	userCount, err := repo.CountUsers(ctx)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"count users failed",
		)
	}

	// 角色分配逻辑：第一个注册用户直接授予管理员权限，其余普通用户
	role := constant.RoleUser
	if userCount == 0 {
		role = constant.RoleAdmin
	}

	// 6. 使用bcrypt对明文密码加密，生成密码哈希存入数据库，不存明文密码
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost, // 默认加密强度
	)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"hash password failed",
		)
	}

	// 7. 组装用户数据库model，准备写入users表
	user := model.User{
		Username:     req.Username,     // 用户名
		Nickname:     req.Nickname,     // 昵称
		Email:        req.Email,        // 邮箱
		PasswordHash: string(passwordHash), // 加密后的密码
		Role:         role,             // 用户角色 admin/user
	}

	// 8. repo层执行数据库插入，创建新用户记录
	if err := repo.CreateUser(ctx, &user); err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"create user failed",
		)
	}

	// 9. 生成全局唯一会话ID，用于本次登录会话标识
	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"generate session failed",
		)
	}

	// 10. 根据新用户信息生成Access+Refresh双令牌，自动处理记住我时效
	tokens, err := utils.GenerateTokenPair(
		user.ID,
		user.Role,
		sessionID,
		req.Remember,
	)
	if err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"generate token failed",
		)
	}

	// 11. 组装会话model，记录登录设备、IP、用户关联信息
	session := model.Session{
		UserID:    user.ID,        // 绑定新注册用户ID
		SessionID: sessionID,      // 会话唯一标识，与JWT载荷对应
		UserIP:    req.UserIP,     // 注册客户端IP，安全审计
		UserAgent: req.UserAgent,  // 浏览器/设备信息
	}

	// 12. 会话写入sessions数据表，完成自动登录
	if err := repo.CreateSession(ctx, &session); err != nil {
		return dto.UserAuthResponse{}, errs.NewInternalServer(
			http.StatusInternalServerError,
			"create session failed",
		)
	}

	// 13. 组装返回给前端的数据：脱敏用户信息 + 双令牌 + 过期时长
	// toUserPrivateResponse：过滤密码等敏感字段，返回安全用户信息
	return dto.UserAuthResponse{
		User:          toUserPrivateResponse(user),
		AccessToken:   tokens.AccessToken,
		RefreshToken:  tokens.RefreshToken,
		AccessMaxAge:  tokens.AccessMaxAge,
		RefreshMaxAge: tokens.RefreshMaxAge,
	}, nil
}

func RequestVerifyEmail(
	ctx context.Context,
	req *dto.VerifyEmailReq,
) error {
	if req == nil {
		return errs.NewBadRequest(
			http.StatusBadRequest,
			"email verification request is required",
		)
	}

	code := utils.RequestEmailVerify(req.Email)
	if code == "" {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"generate email verification code failed",
		)
	}

	body, err := utils.RenderVerificationEmail(code)
	if err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"render verification email failed",
		)
	}

	if err := utils.SendEmail(
		req.Email,
		"邮箱验证码",
		body,
	); err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"send verification email failed",
		)
	}

	return nil
}

func UpdatePassword(
	ctx context.Context,
	userID uint,
	req *dto.UpdatePasswordReq,
) error {
	if userID == 0 || req == nil {
		return errs.NewBadRequest(
			http.StatusBadRequest,
			"invalid password update request",
		)
	}

	user, err := repo.GetUserByID(
		ctx,
		uint64(userID),
	)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.NewUnauthorized(
			http.StatusUnauthorized,
			"user is not authenticated",
		)
	}
	if err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"get user failed",
		)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.OldPassword),
	)
	if err != nil {
		return errs.NewUnauthorized(
			http.StatusUnauthorized,
			"old password is incorrect",
		)
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.NewPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"hash password failed",
		)
	}

	user.PasswordHash = string(passwordHash)

	if err := repo.UpdateUser(ctx, &user); err != nil {
		return errs.NewInternalServer(
			http.StatusInternalServerError,
			"update password failed",
		)
	}

	return nil
}
