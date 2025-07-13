package biz

import "errors"

// 任务相关错误
var (
	ErrUserIDEmpty          = errors.New("userID is required")                // userID不能为空
	ErrInvalidPeriod        = errors.New("invalid period range")              // 时间区间非法
	ErrTitleEmpty           = errors.New("title is required")                 // 标题不能为空
	ErrNoPermission         = errors.New("no permission")                     // 无权操作
	ErrOnlyDayTaskCanScore  = errors.New("only day type task can set score")  // 仅日类型任务可设置分数
	ErrUserIDNotMatchParent = errors.New("userID does not match parent task") // userID与父任务不一致
	ErrTaskNotFound         = errors.New("task not found")                    // 任务不存在
	ErrTaskAlreadyCompleted = errors.New("task already completed")            // 任务已完成
	ErrDuplicateTitle       = errors.New("duplicate title")                   // 标题重复
)

// 日志相关错误
var (
	ErrJournalUserIDEmpty   = errors.New("userID is required")   // userID不能为空
	ErrJournalContentEmpty  = errors.New("content is required")  // 日志内容不能为空
	ErrJournalTypeInvalid   = errors.New("invalid journal type") // 日志类型非法
	ErrJournalPeriodInvalid = errors.New("invalid period range") // 日志时间区间非法
	ErrJournalNotFound      = errors.New("journal not found")    // 日志不存在
)

// 用户相关错误
var (
	ErrUserNameEmpty        = errors.New("username is required")                          // 用户名不能为空
	ErrUserNameDuplicate    = errors.New("username already exists")                       // 用户名已存在
	ErrEmailEmpty           = errors.New("email is required")                             // 邮箱不能为空
	ErrEmailDuplicate       = errors.New("email already exists")                          // 邮箱已存在
	ErrPasswordTooShort     = errors.New("password too short")                            // 密码长度不足
	ErrUserNotFound         = errors.New("user not found")                                // 用户不存在
	ErrUserDeleteNotAllowed = errors.New("user must delete all tasks and journals first") // 用户删除前需先删除所有任务和日志
	ErrPasswordIncorrect    = errors.New("incorrect password")                            // 密码错误
)

// 计划相关错误
var (
	ErrPlanPeriodInvalid = errors.New("invalid plan period")        // 计划时间区间非法
	ErrPlanNoPermission  = errors.New("no permission to view plan") // 无权
)
