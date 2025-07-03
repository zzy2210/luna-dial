package repository

import "errors"

// 定义常用错误
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrTaskNotFound    = errors.New("task not found")
	ErrJournalNotFound = errors.New("journal not found")
)
