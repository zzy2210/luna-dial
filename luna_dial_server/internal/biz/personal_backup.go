package biz

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	// personalBackupSchemaVersion 备份格式版本号
	personalBackupSchemaVersion = 1
	// personalBackupAppName 应用名称（写入 manifest）
	personalBackupAppName = "luna-dial-server"
	// personalBackupAppVersion 应用版本（写入 manifest）
	personalBackupAppVersion = "1.0.0"
	// personalBackupManifestFilename 清单文件名
	personalBackupManifestFilename = "manifest.json"
	// personalBackupTasksFilename 任务数据文件名（JSONL）
	personalBackupTasksFilename = "tasks.jsonl"
	// personalBackupJournalsFilename 日志数据文件名（JSONL）
	personalBackupJournalsFilename = "journals.jsonl"
)

// PersonalBackupMaxBytes 个人备份上传大小上限（用于服务端限流/保护）
const PersonalBackupMaxBytes int64 = 50 << 20 // 50MB

// PersonalBackupError 个人备份相关错误（用于向 HTTP 层返回可控状态码）
type PersonalBackupError struct {
	// status HTTP 状态码
	status int
	// message 错误信息
	message string
}

func (e *PersonalBackupError) Error() string {
	return e.message
}

// HTTPStatus 返回 HTTP 状态码
func (e *PersonalBackupError) HTTPStatus() int {
	return e.status
}

func newPersonalBackupError(status int, format string, args ...any) error {
	return &PersonalBackupError{
		status:  status,
		message: fmt.Sprintf(format, args...),
	}
}

// PersonalBackupImportResult 个人备份导入结果
type PersonalBackupImportResult struct {
	// ImportedTasks 导入的任务数量
	ImportedTasks int `json:"imported_tasks"`
	// ImportedJournals 导入的日志数量
	ImportedJournals int `json:"imported_journals"`
}

// PersonalBackupUsecase 个人备份用例层
type PersonalBackupUsecase struct {
	// repo 个人备份仓库
	repo PersonalBackupRepo
}

// NewPersonalBackupUsecase 创建个人备份用例层
func NewPersonalBackupUsecase(repo PersonalBackupRepo) *PersonalBackupUsecase {
	return &PersonalBackupUsecase{repo: repo}
}

// ExportArchive 导出当前用户的个人备份包（tar.gz）
func (uc *PersonalBackupUsecase) ExportArchive(ctx context.Context, userID string, username string, now time.Time) ([]byte, error) {
	if userID == "" {
		return nil, newPersonalBackupError(401, "missing user_id")
	}
	if username == "" {
		return nil, newPersonalBackupError(401, "missing username")
	}

	tasks, err := uc.repo.ListAllTasksForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	journals, err := uc.repo.ListAllJournalsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	tasksBytes, err := encodeTasksJSONL(tasks)
	if err != nil {
		return nil, err
	}
	journalsBytes, err := encodeJournalsJSONL(journals)
	if err != nil {
		return nil, err
	}

	manifest := personalBackupManifest{
		SchemaVersion: personalBackupSchemaVersion,
		CreatedAt:     now.UTC(),
		Subject: personalBackupManifestSubject{
			Username: username,
		},
		App: personalBackupManifestApp{
			Name:    personalBackupAppName,
			Version: personalBackupAppVersion,
		},
		Content: personalBackupManifestContent{
			Tasks: personalBackupManifestContentItem{
				Filename: personalBackupTasksFilename,
				SHA256:   sha256Hex(tasksBytes),
				Count:    len(tasks),
			},
			Journals: personalBackupManifestContentItem{
				Filename: personalBackupJournalsFilename,
				SHA256:   sha256Hex(journalsBytes),
				Count:    len(journals),
			},
		},
	}

	manifestBytes, err := json.MarshalIndent(&manifest, "", "  ")
	if err != nil {
		return nil, err
	}

	archiveBytes, err := buildTarGz(map[string][]byte{
		personalBackupManifestFilename: manifestBytes,
		personalBackupTasksFilename:    tasksBytes,
		personalBackupJournalsFilename: journalsBytes,
	})
	if err != nil {
		return nil, err
	}

	return archiveBytes, nil
}

// ImportOverwrite 覆盖式导入备份包
func (uc *PersonalBackupUsecase) ImportOverwrite(ctx context.Context, userID string, username string, force bool, archiveBytes []byte) (*PersonalBackupImportResult, error) {
	if userID == "" {
		return nil, newPersonalBackupError(401, "missing user_id")
	}
	if username == "" {
		return nil, newPersonalBackupError(401, "missing username")
	}
	if len(archiveBytes) == 0 {
		return nil, newPersonalBackupError(400, "empty archive")
	}

	files, err := readTarGzFiles(bytes.NewReader(archiveBytes))
	if err != nil {
		return nil, newPersonalBackupError(422, "invalid archive: %v", err)
	}

	manifestBytes, ok := files[personalBackupManifestFilename]
	if !ok {
		return nil, newPersonalBackupError(422, "missing %s", personalBackupManifestFilename)
	}

	// personalBackupManifest 备份清单结构
	var manifest personalBackupManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, newPersonalBackupError(422, "invalid manifest.json: %v", err)
	}
	if manifest.SchemaVersion != personalBackupSchemaVersion {
		return nil, newPersonalBackupError(422, "unsupported schema_version: %d", manifest.SchemaVersion)
	}

	if !force && manifest.Subject.Username != "" && manifest.Subject.Username != username {
		return nil, newPersonalBackupError(422, "backup subject username mismatch: %q != %q (use force=true to override)", manifest.Subject.Username, username)
	}

	tasksBytes, ok := files[personalBackupTasksFilename]
	if !ok {
		return nil, newPersonalBackupError(422, "missing %s", personalBackupTasksFilename)
	}
	journalsBytes, ok := files[personalBackupJournalsFilename]
	if !ok {
		return nil, newPersonalBackupError(422, "missing %s", personalBackupJournalsFilename)
	}

	if want := strings.TrimSpace(manifest.Content.Tasks.Filename); want != "" && want != personalBackupTasksFilename {
		return nil, newPersonalBackupError(422, "tasks filename mismatch")
	}
	if want := strings.TrimSpace(manifest.Content.Journals.Filename); want != "" && want != personalBackupJournalsFilename {
		return nil, newPersonalBackupError(422, "journals filename mismatch")
	}

	// sha256 校验（允许 manifest 未填字段时跳过，但一旦填了就必须匹配）
	if want := strings.TrimSpace(manifest.Content.Tasks.SHA256); want != "" {
		got := sha256Hex(tasksBytes)
		if !strings.EqualFold(want, got) {
			return nil, newPersonalBackupError(422, "tasks sha256 mismatch")
		}
	}
	if want := strings.TrimSpace(manifest.Content.Journals.SHA256); want != "" {
		got := sha256Hex(journalsBytes)
		if !strings.EqualFold(want, got) {
			return nil, newPersonalBackupError(422, "journals sha256 mismatch")
		}
	}

	tasks, err := decodeTasksJSONL(tasksBytes)
	if err != nil {
		return nil, err
	}
	journals, err := decodeJournalsJSONL(journalsBytes)
	if err != nil {
		return nil, err
	}

	if manifest.Content.Tasks.Count != len(tasks) {
		return nil, newPersonalBackupError(422, "tasks count mismatch")
	}
	if manifest.Content.Journals.Count != len(journals) {
		return nil, newPersonalBackupError(422, "journals count mismatch")
	}

	// 防止“坏备份”把树字段搞崩：导入前统一重建 derived fields（中文名：派生冗余字段）
	if err := rebuildTaskTreeFields(tasks); err != nil {
		return nil, err
	}

	if err := uc.repo.ImportOverwrite(ctx, userID, tasks, journals); err != nil {
		return nil, err
	}

	return &PersonalBackupImportResult{
		ImportedTasks:    len(tasks),
		ImportedJournals: len(journals),
	}, nil
}

// personalBackupManifest 备份清单结构
type personalBackupManifest struct {
	// SchemaVersion 备份格式版本号
	SchemaVersion int `json:"schema_version"`
	// CreatedAt 备份创建时间
	CreatedAt time.Time `json:"created_at"`
	// Subject 备份主体信息
	Subject personalBackupManifestSubject `json:"subject"`
	// App 导出应用信息
	App personalBackupManifestApp `json:"app"`
	// Content 内容校验与统计
	Content personalBackupManifestContent `json:"content"`
}

// personalBackupManifestSubject 备份主体信息
type personalBackupManifestSubject struct {
	// Username 备份主体用户名
	Username string `json:"username"`
}

// personalBackupManifestApp 应用信息
type personalBackupManifestApp struct {
	// Name 应用名称
	Name string `json:"name"`
	// Version 应用版本
	Version string `json:"version"`
}

// personalBackupManifestContentItem 单个文件条目
type personalBackupManifestContentItem struct {
	// Filename 文件名
	Filename string `json:"filename"`
	// SHA256 文件 sha256
	SHA256 string `json:"sha256"`
	// Count 记录条数
	Count int `json:"count"`
}

// personalBackupManifestContent 内容清单
type personalBackupManifestContent struct {
	// Tasks 任务条目
	Tasks personalBackupManifestContentItem `json:"tasks"`
	// Journals 日志条目
	Journals personalBackupManifestContentItem `json:"journals"`
}

// personalBackupTaskRecord 任务备份记录（JSONL 单行）
type personalBackupTaskRecord struct {
	// ID 任务ID
	ID string `json:"id"`
	// Title 标题
	Title string `json:"title"`
	// TaskType 周期类型
	TaskType string `json:"task_type"`
	// PeriodStart 周期起始
	PeriodStart *time.Time `json:"period_start"`
	// PeriodEnd 周期结束
	PeriodEnd *time.Time `json:"period_end"`
	// Tags 标签
	Tags []string `json:"tags"`
	// Icon 图标
	Icon string `json:"icon"`
	// Score 分数
	Score int `json:"score"`
	// Status 状态
	Status string `json:"status"`
	// Priority 优先级
	Priority string `json:"priority"`
	// ParentID 父任务ID
	ParentID string `json:"parent_id"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

// personalBackupJournalRecord 日志备份记录（JSONL 单行）
type personalBackupJournalRecord struct {
	// ID 日志ID
	ID string `json:"id"`
	// Title 标题
	Title string `json:"title"`
	// Content 内容
	Content string `json:"content"`
	// JournalType 周期类型
	JournalType string `json:"journal_type"`
	// PeriodStart 周期起始
	PeriodStart *time.Time `json:"period_start"`
	// PeriodEnd 周期结束
	PeriodEnd *time.Time `json:"period_end"`
	// Icon 图标
	Icon string `json:"icon"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at"`
}

func encodeTasksJSONL(tasks []*Task) ([]byte, error) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	for i := range tasks {
		if tasks[i] == nil {
			continue
		}
		rec, err := taskToBackupRecord(tasks[i])
		if err != nil {
			return nil, err
		}
		line, err := json.Marshal(rec)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(line); err != nil {
			return nil, err
		}
		if err := w.WriteByte('\n'); err != nil {
			return nil, err
		}
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeJournalsJSONL(journals []*Journal) ([]byte, error) {
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)

	for i := range journals {
		if journals[i] == nil {
			continue
		}
		rec, err := journalToBackupRecord(journals[i])
		if err != nil {
			return nil, err
		}
		line, err := json.Marshal(rec)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(line); err != nil {
			return nil, err
		}
		if err := w.WriteByte('\n'); err != nil {
			return nil, err
		}
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeTasksJSONL(in []byte) ([]*Task, error) {
	recs, err := decodeJSONLLines[personalBackupTaskRecord](in)
	if err != nil {
		return nil, err
	}

	out := make([]*Task, 0, len(recs))
	seen := make(map[string]struct{}, len(recs))
	for i := range recs {
		rec := recs[i]
		if strings.TrimSpace(rec.ID) == "" {
			return nil, newPersonalBackupError(422, "task id is required (line %d)", i+1)
		}
		if _, ok := seen[rec.ID]; ok {
			return nil, newPersonalBackupError(422, "duplicate task id: %s", rec.ID)
		}
		seen[rec.ID] = struct{}{}

		task, err := backupRecordToTask(&rec)
		if err != nil {
			return nil, newPersonalBackupError(422, "invalid task record (line %d): %v", i+1, err)
		}
		out = append(out, task)
	}

	return out, nil
}

func decodeJournalsJSONL(in []byte) ([]*Journal, error) {
	recs, err := decodeJSONLLines[personalBackupJournalRecord](in)
	if err != nil {
		return nil, err
	}

	out := make([]*Journal, 0, len(recs))
	seen := make(map[string]struct{}, len(recs))
	for i := range recs {
		rec := recs[i]
		if strings.TrimSpace(rec.ID) == "" {
			return nil, newPersonalBackupError(422, "journal id is required (line %d)", i+1)
		}
		if _, ok := seen[rec.ID]; ok {
			return nil, newPersonalBackupError(422, "duplicate journal id: %s", rec.ID)
		}
		seen[rec.ID] = struct{}{}

		journal, err := backupRecordToJournal(&rec)
		if err != nil {
			return nil, newPersonalBackupError(422, "invalid journal record (line %d): %v", i+1, err)
		}
		out = append(out, journal)
	}

	return out, nil
}

// personalBackupLineDecodeError JSONL 解析错误
type personalBackupLineDecodeError struct {
	// Line 行号（1-based）
	Line int
	// Err 原始错误
	Err error
}

func (e *personalBackupLineDecodeError) Error() string {
	return fmt.Sprintf("line %d: %v", e.Line, e.Err)
}

func decodeJSONLLines[T any](in []byte) ([]T, error) {
	r := bufio.NewReader(bytes.NewReader(in))
	out := make([]T, 0, 128)

	lineNo := 0
	for {
		b, err := r.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, newPersonalBackupError(422, "invalid jsonl: %v", err)
		}
		if len(b) == 0 && errors.Is(err, io.EOF) {
			break
		}

		lineNo++
		line := bytes.TrimSpace(b)
		if len(line) == 0 {
			if errors.Is(err, io.EOF) {
				break
			}
			continue
		}

		var rec T
		if jerr := json.Unmarshal(line, &rec); jerr != nil {
			return nil, newPersonalBackupError(422, "invalid json (line %d)", lineNo)
		}
		out = append(out, rec)

		if errors.Is(err, io.EOF) {
			break
		}
	}

	return out, nil
}

func taskToBackupRecord(t *Task) (*personalBackupTaskRecord, error) {
	if t.CreatedAt.IsZero() {
		return nil, newPersonalBackupError(500, "task created_at is missing: %s", t.ID)
	}
	if t.UpdatedAt.IsZero() {
		return nil, newPersonalBackupError(500, "task updated_at is missing: %s", t.ID)
	}

	taskType, err := periodTypeToString(t.TaskType)
	if err != nil {
		return nil, err
	}
	status, err := taskStatusToString(t.Status)
	if err != nil {
		return nil, err
	}
	priority, err := taskPriorityToString(t.Priority)
	if err != nil {
		return nil, err
	}

	start := timePtrOrNil(t.TimePeriod.Start)
	end := timePtrOrNil(t.TimePeriod.End)

	return &personalBackupTaskRecord{
		ID:          t.ID,
		Title:       t.Title,
		TaskType:    taskType,
		PeriodStart: start,
		PeriodEnd:   end,
		Tags:        append([]string(nil), t.Tags...),
		Icon:        t.Icon,
		Score:       t.Score,
		Status:      status,
		Priority:    priority,
		ParentID:    t.ParentID,
		CreatedAt:   t.CreatedAt.UTC(),
		UpdatedAt:   t.UpdatedAt.UTC(),
	}, nil
}

func journalToBackupRecord(j *Journal) (*personalBackupJournalRecord, error) {
	if j.CreatedAt.IsZero() {
		return nil, newPersonalBackupError(500, "journal created_at is missing: %s", j.ID)
	}
	if j.UpdatedAt.IsZero() {
		return nil, newPersonalBackupError(500, "journal updated_at is missing: %s", j.ID)
	}

	jType, err := periodTypeToString(j.JournalType)
	if err != nil {
		return nil, err
	}

	start := timePtrOrNil(j.TimePeriod.Start)
	end := timePtrOrNil(j.TimePeriod.End)

	return &personalBackupJournalRecord{
		ID:          j.ID,
		Title:       j.Title,
		Content:     j.Content,
		JournalType: jType,
		PeriodStart: start,
		PeriodEnd:   end,
		Icon:        j.Icon,
		CreatedAt:   j.CreatedAt.UTC(),
		UpdatedAt:   j.UpdatedAt.UTC(),
	}, nil
}

func backupRecordToTask(r *personalBackupTaskRecord) (*Task, error) {
	pt, err := periodTypeFromString(r.TaskType)
	if err != nil {
		return nil, err
	}
	status, err := taskStatusFromString(r.Status)
	if err != nil {
		return nil, err
	}
	priority, err := taskPriorityFromString(r.Priority)
	if err != nil {
		return nil, err
	}
	if r.CreatedAt.IsZero() {
		return nil, fmt.Errorf("created_at is required")
	}
	if r.UpdatedAt.IsZero() {
		return nil, fmt.Errorf("updated_at is required")
	}

	start := time.Time{}
	end := time.Time{}
	if r.PeriodStart != nil {
		start = r.PeriodStart.UTC()
	}
	if r.PeriodEnd != nil {
		end = r.PeriodEnd.UTC()
	}

	return &Task{
		ID:       r.ID,
		Title:    r.Title,
		TaskType: pt,
		TimePeriod: Period{
			Start: start,
			End:   end,
		},
		Tags:          append([]string(nil), r.Tags...),
		Icon:          r.Icon,
		Score:         r.Score,
		Status:        status,
		Priority:      priority,
		ParentID:      r.ParentID,
		UserID:        "",
		CreatedAt:     r.CreatedAt.UTC(),
		UpdatedAt:     r.UpdatedAt.UTC(),
		HasChildren:   false,
		ChildrenCount: 0,
		RootTaskID:    "",
		TreeDepth:     0,
		Children:      make([]*Task, 0),
	}, nil
}

func backupRecordToJournal(r *personalBackupJournalRecord) (*Journal, error) {
	pt, err := periodTypeFromString(r.JournalType)
	if err != nil {
		return nil, err
	}
	if r.CreatedAt.IsZero() {
		return nil, fmt.Errorf("created_at is required")
	}
	if r.UpdatedAt.IsZero() {
		return nil, fmt.Errorf("updated_at is required")
	}

	start := time.Time{}
	end := time.Time{}
	if r.PeriodStart != nil {
		start = r.PeriodStart.UTC()
	}
	if r.PeriodEnd != nil {
		end = r.PeriodEnd.UTC()
	}

	return &Journal{
		ID:          r.ID,
		Title:       r.Title,
		Content:     r.Content,
		JournalType: pt,
		TimePeriod: Period{
			Start: start,
			End:   end,
		},
		Icon:      r.Icon,
		UserID:    "",
		CreatedAt: r.CreatedAt.UTC(),
		UpdatedAt: r.UpdatedAt.UTC(),
	}, nil
}

func buildTarGz(files map[string][]byte) ([]byte, error) {
	var out bytes.Buffer
	gw := gzip.NewWriter(&out)
	tw := tar.NewWriter(gw)

	writeFile := func(name string, b []byte) error {
		hdr := &tar.Header{
			Name:    name,
			Mode:    0o644,
			Size:    int64(len(b)),
			ModTime: time.Now().UTC(),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write(b); err != nil {
			return err
		}
		return nil
	}

	// 固定顺序，避免不同运行产生不必要的差异
	if err := writeFile(personalBackupManifestFilename, files[personalBackupManifestFilename]); err != nil {
		return nil, err
	}
	if err := writeFile(personalBackupTasksFilename, files[personalBackupTasksFilename]); err != nil {
		return nil, err
	}
	if err := writeFile(personalBackupJournalsFilename, files[personalBackupJournalsFilename]); err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func readTarGzFiles(r io.Reader) (map[string][]byte, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	out := make(map[string][]byte, 8)

	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr == nil {
			continue
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		if hdr.Size < 0 {
			return nil, fmt.Errorf("invalid file size: %s", hdr.Name)
		}

		b, err := io.ReadAll(io.LimitReader(tr, hdr.Size))
		if err != nil {
			return nil, err
		}
		out[hdr.Name] = b
	}

	return out, nil
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func timePtrOrNil(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	tt := t.UTC()
	return &tt
}

func periodTypeToString(pt PeriodType) (string, error) {
	switch pt {
	case PeriodDay:
		return "day", nil
	case PeriodWeek:
		return "week", nil
	case PeriodMonth:
		return "month", nil
	case PeriodQuarter:
		return "quarter", nil
	case PeriodYear:
		return "year", nil
	default:
		return "", newPersonalBackupError(422, "unknown period type: %d", int(pt))
	}
}

func periodTypeFromString(s string) (PeriodType, error) {
	switch strings.TrimSpace(s) {
	case "day":
		return PeriodDay, nil
	case "week":
		return PeriodWeek, nil
	case "month":
		return PeriodMonth, nil
	case "quarter":
		return PeriodQuarter, nil
	case "year":
		return PeriodYear, nil
	default:
		return 0, fmt.Errorf("unknown period type: %s", s)
	}
}

func taskStatusToString(st TaskStatus) (string, error) {
	switch st {
	case TaskStatusNotStarted:
		return "not_started", nil
	case TaskStatusInProgress:
		return "in_progress", nil
	case TaskStatusCompleted:
		return "completed", nil
	case TaskStatusCancelled:
		return "cancelled", nil
	default:
		return "", newPersonalBackupError(422, "unknown task status: %d", int(st))
	}
}

func taskStatusFromString(s string) (TaskStatus, error) {
	switch strings.TrimSpace(s) {
	case "not_started":
		return TaskStatusNotStarted, nil
	case "in_progress":
		return TaskStatusInProgress, nil
	case "completed":
		return TaskStatusCompleted, nil
	case "cancelled":
		return TaskStatusCancelled, nil
	default:
		return 0, fmt.Errorf("unknown task status: %s", s)
	}
}

func taskPriorityToString(p TaskPriority) (string, error) {
	switch p {
	case TaskPriorityLow:
		return "low", nil
	case TaskPriorityMedium:
		return "medium", nil
	case TaskPriorityHigh:
		return "high", nil
	case TaskPriorityUrgent:
		return "urgent", nil
	default:
		return "", newPersonalBackupError(422, "unknown task priority: %d", int(p))
	}
}

func taskPriorityFromString(s string) (TaskPriority, error) {
	switch strings.TrimSpace(s) {
	case "low":
		return TaskPriorityLow, nil
	case "medium":
		return TaskPriorityMedium, nil
	case "high":
		return TaskPriorityHigh, nil
	case "urgent":
		return TaskPriorityUrgent, nil
	default:
		return 0, fmt.Errorf("unknown task priority: %s", s)
	}
}

func rebuildTaskTreeFields(tasks []*Task) error {
	if len(tasks) == 0 {
		return nil
	}

	byID := make(map[string]*Task, len(tasks))
	for i := range tasks {
		t := tasks[i]
		if t == nil {
			continue
		}
		// 重置 derived fields（中文名：派生冗余字段）
		t.HasChildren = false
		t.ChildrenCount = 0
		t.RootTaskID = t.ID
		t.TreeDepth = 0
		byID[t.ID] = t
	}

	// 先统计直接子任务数量
	for _, t := range byID {
		if t.ParentID == "" {
			continue
		}
		parent, ok := byID[t.ParentID]
		if !ok {
			return newPersonalBackupError(422, "task parent_id not found: %s -> %s", t.ID, t.ParentID)
		}
		parent.ChildrenCount++
		parent.HasChildren = true
	}

	// 再计算 root_task_id / tree_depth，带环检测 + 记忆化避免 O(n^2)
	type state struct {
		// visiting DFS 中访问中标记
		visiting bool
		// done 计算完成标记
		done bool
	}
	states := make(map[string]*state, len(byID))

	var walk func(t *Task) (string, int, error)
	walk = func(t *Task) (string, int, error) {
		if t == nil {
			return "", 0, fmt.Errorf("nil task")
		}
		s := states[t.ID]
		if s == nil {
			s = &state{}
			states[t.ID] = s
		}
		if s.done {
			return t.RootTaskID, t.TreeDepth, nil
		}
		if s.visiting {
			return "", 0, newPersonalBackupError(422, "task parent cycle detected at %s", t.ID)
		}
		s.visiting = true

		if t.ParentID == "" {
			t.RootTaskID = t.ID
			t.TreeDepth = 0
		} else {
			parent := byID[t.ParentID]
			if parent == nil {
				return "", 0, newPersonalBackupError(422, "task parent_id not found: %s -> %s", t.ID, t.ParentID)
			}
			rootID, depth, err := walk(parent)
			if err != nil {
				return "", 0, err
			}
			t.RootTaskID = rootID
			t.TreeDepth = depth + 1
		}

		s.visiting = false
		s.done = true
		return t.RootTaskID, t.TreeDepth, nil
	}

	for _, t := range byID {
		if _, _, err := walk(t); err != nil {
			return err
		}
	}

	return nil
}
