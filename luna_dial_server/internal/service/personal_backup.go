package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"luna_dial/internal/biz"
)

// handlePersonalBackupExport 导出当前用户的个人数据备份包（tar.gz）
func (s *Service) handlePersonalBackupExport(c echo.Context) error {
	userID, username, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User information not found"))
	}

	now := time.Now().UTC()
	archiveBytes, err := s.personalBackupUsecase.ExportArchive(c.Request().Context(), userID, username, now)
	if err != nil {
		var pbErr *biz.PersonalBackupError
		if errors.As(err, &pbErr) {
			return c.JSON(pbErr.HTTPStatus(), NewErrorResponse(pbErr.HTTPStatus(), pbErr.Error()))
		}
		return c.JSON(500, NewErrorResponse(500, "Failed to build backup archive"))
	}

	filename := fmt.Sprintf("luna-dial-personal-backup-%s.tar.gz", now.Format("20060102-150405"))
	c.Response().Header().Set(echo.HeaderContentType, "application/gzip")
	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%q", filename))
	c.Response().WriteHeader(http.StatusOK)
	if _, err := c.Response().Writer.Write(archiveBytes); err != nil {
		// 响应可能已部分写出，尽量避免二次写 JSON；直接返回错误让 Echo 记录日志。
		return err
	}
	return nil
}

// handlePersonalBackupImportOverwrite 导入当前用户的个人数据备份包（覆盖式）
func (s *Service) handlePersonalBackupImportOverwrite(c echo.Context) error {
	force := c.QueryParam("force") == "true"

	userID, username, err := GetUserFromContext(c)
	if err != nil {
		return c.JSON(401, NewErrorResponse(401, "User information not found"))
	}

	// 限制请求体大小，避免 multipart 在解析时读入过大内容。
	req := c.Request()
	req.Body = http.MaxBytesReader(c.Response().Writer, req.Body, biz.PersonalBackupMaxBytes)

	formFile, err := c.FormFile("file")
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, "file is required"))
	}

	f, err := formFile.Open()
	if err != nil {
		return c.JSON(400, NewErrorResponse(400, "failed to open uploaded file"))
	}
	defer f.Close()

	archiveBytes, err := io.ReadAll(f)
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return c.JSON(413, NewErrorResponse(413, "backup file is too large"))
		}
		return c.JSON(400, NewErrorResponse(400, "failed to read uploaded file"))
	}

	res, err := s.personalBackupUsecase.ImportOverwrite(c.Request().Context(), userID, username, force, archiveBytes)
	if err != nil {
		var pbErr *biz.PersonalBackupError
		if errors.As(err, &pbErr) {
			return c.JSON(pbErr.HTTPStatus(), NewErrorResponse(pbErr.HTTPStatus(), pbErr.Error()))
		}
		return c.JSON(500, NewErrorResponse(500, "failed to import backup"))
	}

	return c.JSON(200, NewSuccessResponse(map[string]interface{}{
		"imported_tasks":    res.ImportedTasks,
		"imported_journals": res.ImportedJournals,
	}))
}
