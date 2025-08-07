package data

import (
	"luna_dial/internal/biz"
	"strings"
)

// TaskConverter 任务数据转换器
type TaskConverter struct{}

func NewTaskConverter() *TaskConverter {
	return &TaskConverter{}
}

// BizToData 业务模型转数据模型
func (c *TaskConverter) BizToData(bizTask *biz.Task) *Task {
	if bizTask == nil {
		return nil
	}

	dataTask := &Task{
		ID:          bizTask.ID,
		UserID:      bizTask.UserID,
		Title:       bizTask.Title,
		TaskType:    int(bizTask.TaskType),
		PeriodStart: bizTask.TimePeriod.Start,
		PeriodEnd:   bizTask.TimePeriod.End,
		Score:       bizTask.Score,
		Status:      int(bizTask.Status),
		Priority:    int(bizTask.Priority),
		ParentID:    bizTask.ParentID,
		Icon:        bizTask.Icon,
		CreatedAt:   bizTask.CreatedAt,
		UpdatedAt:   bizTask.UpdatedAt,
		
		// 新增：树结构优化字段转换
		// 这些字段直接从业务层同步到数据层，确保数据一致性
		HasChildren:   bizTask.HasChildren,
		ChildrenCount: bizTask.ChildrenCount,
		RootTaskID:    bizTask.RootTaskID,
		TreeDepth:     bizTask.TreeDepth,
	}

	// 处理Tags数组转逗号分隔字符串
	if len(bizTask.Tags) > 0 {
		dataTask.Tags = strings.Join(bizTask.Tags, ",")
	}

	return dataTask
}

// DataToBiz 数据模型转业务模型
func (c *TaskConverter) DataToBiz(dataTask *Task) *biz.Task {
	if dataTask == nil {
		return nil
	}

	bizTask := &biz.Task{
		ID:       dataTask.ID,
		UserID:   dataTask.UserID,
		Title:    dataTask.Title,
		TaskType: biz.PeriodType(dataTask.TaskType),
		TimePeriod: biz.Period{
			Start: dataTask.PeriodStart,
			End:   dataTask.PeriodEnd,
		},
		Score:    dataTask.Score,
		Status:   biz.TaskStatus(dataTask.Status),
		Priority: biz.TaskPriority(dataTask.Priority),
		ParentID: dataTask.ParentID,
		Icon:        dataTask.Icon,
		CreatedAt:   dataTask.CreatedAt,
		UpdatedAt:   dataTask.UpdatedAt,
		
		// 新增：树结构优化字段转换
		// 从数据库字段同步到业务层，为后续树构建提供基础数据
		HasChildren:   dataTask.HasChildren,
		ChildrenCount: dataTask.ChildrenCount,
		RootTaskID:    dataTask.RootTaskID,
		TreeDepth:     dataTask.TreeDepth,
		
		// Children字段在这里初始化为空切片，由上层业务逻辑负责构建树结构
		// 设计思路：转换器只负责基础数据转换，树关系构建由专门的业务方法处理
		Children: make([]*biz.Task, 0),
	}

	// 处理Tags逗号分隔字符串转数组
	if dataTask.Tags != "" {
		bizTask.Tags = strings.Split(dataTask.Tags, ",")
		// 去除空字符串元素
		validTags := make([]string, 0, len(bizTask.Tags))
		for _, tag := range bizTask.Tags {
			if trimmed := strings.TrimSpace(tag); trimmed != "" {
				validTags = append(validTags, trimmed)
			}
		}
		bizTask.Tags = validTags
	}

	return bizTask
}

// DataToBizList 批量数据模型转业务模型
func (c *TaskConverter) DataToBizList(dataTasks []*Task) []*biz.Task {
	if len(dataTasks) == 0 {
		return nil
	}

	bizTasks := make([]*biz.Task, len(dataTasks))
	for i, dataTask := range dataTasks {
		bizTasks[i] = c.DataToBiz(dataTask)
	}
	return bizTasks
}

// BizTDataList 批量业务模型转数据模型
func (c *TaskConverter) BizToDataList(bizTasks []*biz.Task) []*Task {
	if len(bizTasks) == 0 {
		return nil
	}

	dataTasks := make([]*Task, len(bizTasks))
	for i, bizTask := range bizTasks {
		dataTasks[i] = c.BizToData(bizTask)
	}
	return dataTasks
}

// JournalConverter 日志数据转换器
type JournalConverter struct{}

func NewJournalConverter() *JournalConverter {
	return &JournalConverter{}
}

// BizToData 业务模型转数据模型
func (c *JournalConverter) BizToData(bizJournal *biz.Journal) *Journal {
	if bizJournal == nil {
		return nil
	}

	return &Journal{
		ID:          bizJournal.ID,
		UserID:      bizJournal.UserID,
		Title:       bizJournal.Title,
		Content:     bizJournal.Content,
		JournalType: int(bizJournal.JournalType),
		PeriodStart: bizJournal.TimePeriod.Start,
		PeriodEnd:   bizJournal.TimePeriod.End,
		Icon:        bizJournal.Icon,
		CreatedAt:   bizJournal.CreatedAt,
		UpdatedAt:   bizJournal.UpdatedAt,
	}
}

// DataToBiz 数据模型转业务模型
func (c *JournalConverter) DataToBiz(dataJournal *Journal) *biz.Journal {
	if dataJournal == nil {
		return nil
	}

	return &biz.Journal{
		ID:          dataJournal.ID,
		UserID:      dataJournal.UserID,
		Title:       dataJournal.Title,
		Content:     dataJournal.Content,
		JournalType: biz.PeriodType(dataJournal.JournalType),
		TimePeriod: biz.Period{
			Start: dataJournal.PeriodStart,
			End:   dataJournal.PeriodEnd,
		},
		Icon:      dataJournal.Icon,
		CreatedAt: dataJournal.CreatedAt,
		UpdatedAt: dataJournal.UpdatedAt,
	}
}

// DataToBizList 批量转换
func (c *JournalConverter) DataToBizList(dataJournals []*Journal) []*biz.Journal {
	if len(dataJournals) == 0 {
		return nil
	}

	bizJournals := make([]*biz.Journal, len(dataJournals))
	for i, dataJournal := range dataJournals {
		bizJournals[i] = c.DataToBiz(dataJournal)
	}
	return bizJournals
}

// BizToDataList 批量业务模型转数据模型
func (c *JournalConverter) BizToDataList(bizJournals []*biz.Journal) []*Journal {
	if len(bizJournals) == 0 {
		return nil
	}

	dataJournals := make([]*Journal, len(bizJournals))
	for i, bizJournal := range bizJournals {
		dataJournals[i] = c.BizToData(bizJournal)
	}
	return dataJournals
}

// UserConverter 用户数据转换器
type UserConverter struct{}

func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// BizToData 业务模型转数据模型
func (c *UserConverter) BizToData(bizUser *biz.User) *User {
	if bizUser == nil {
		return nil
	}

	return &User{
		ID:        bizUser.ID,
		UserName:  bizUser.Username,
		Name:      bizUser.Name,
		Email:     bizUser.Email,
		Password:  bizUser.Password,
		CreatedAt: bizUser.CreatedAt,
		UpdatedAt: bizUser.UpdatedAt,
	}
}

// DataToBiz 数据模型转业务模型
func (c *UserConverter) DataToBiz(dataUser *User) *biz.User {
	if dataUser == nil {
		return nil
	}

	return &biz.User{
		ID:        dataUser.ID,
		Username:  dataUser.UserName,
		Name:      dataUser.Name,
		Email:     dataUser.Email,
		Password:  dataUser.Password,
		CreatedAt: dataUser.CreatedAt,
		UpdatedAt: dataUser.UpdatedAt,
	}
}

// DataToBizList 批量数据模型转业务模型
func (c *UserConverter) DataToBizList(dataUsers []*User) []*biz.User {
	if len(dataUsers) == 0 {
		return nil
	}

	bizUsers := make([]*biz.User, len(dataUsers))
	for i, dataUser := range dataUsers {
		bizUsers[i] = c.DataToBiz(dataUser)
	}
	return bizUsers
}

// BizToDataList 批量业务模型转数据模型
func (c *UserConverter) BizToDataList(bizUsers []*biz.User) []*User {
	if len(bizUsers) == 0 {
		return nil
	}

	dataUsers := make([]*User, len(bizUsers))
	for i, bizUser := range bizUsers {
		dataUsers[i] = c.BizToData(bizUser)
	}
	return dataUsers
}
