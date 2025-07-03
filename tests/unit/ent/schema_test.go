package ent_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"okr-web/ent/enttest"
	"okr-web/ent/task"
	"okr-web/ent/user"
)

var ctx = context.Background()

func TestSchemaBasics(t *testing.T) {
	// 使用内存SQLite数据库进行测试
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// 验证能够创建User
	createdUser, err := client.User.Create().
		SetUsername("testuser").
		SetPassword("hashedpassword").
		SetEmail("test@example.com").
		Save(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, createdUser.ID)
	assert.Equal(t, "testuser", createdUser.Username)
	assert.Equal(t, "test@example.com", createdUser.Email)

	// 验证能够创建Task
	createdTask, err := client.Task.Create().
		SetTitle("测试任务").
		SetDescription("这是一个测试任务").
		SetType(task.TypeDay).
		SetStatus(task.StatusPending).
		SetScore(5).
		SetUser(createdUser).
		Save(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, createdTask.ID)
	assert.Equal(t, "测试任务", createdTask.Title)
	assert.Equal(t, task.TypeDay, createdTask.Type)
	assert.Equal(t, task.StatusPending, createdTask.Status)
	assert.Equal(t, 5, createdTask.Score)

	// 验证能够查询数据
	foundUser, err := client.User.Query().
		Where(user.UsernameEQ("testuser")).
		Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, foundUser.ID)

	// 验证关联查询
	userWithTasks, err := client.User.Query().
		Where(user.IDEQ(createdUser.ID)).
		WithTasks().
		Only(ctx)
	require.NoError(t, err)
	assert.Len(t, userWithTasks.Edges.Tasks, 1)
	assert.Equal(t, createdTask.ID, userWithTasks.Edges.Tasks[0].ID)

	t.Log("✅ Ent schema generation and basic operations working correctly")
}

func TestTaskHierarchy(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	// 创建用户
	testUser, err := client.User.Create().
		SetUsername("hierarchyuser").
		SetPassword("password").
		Save(ctx)
	require.NoError(t, err)

	// 创建父任务
	parentTask, err := client.Task.Create().
		SetTitle("父任务").
		SetType(task.TypeYear).
		SetUser(testUser).
		Save(ctx)
	require.NoError(t, err)

	// 创建子任务
	childTask, err := client.Task.Create().
		SetTitle("子任务").
		SetType(task.TypeQuarter).
		SetUser(testUser).
		SetParent(parentTask).
		Save(ctx)
	require.NoError(t, err)

	// 验证层级关系
	assert.Equal(t, parentTask.ID, *childTask.ParentID)

	// 验证查询父子关系
	parentWithChildren, err := client.Task.Query().
		Where(task.IDEQ(parentTask.ID)).
		WithChildren().
		Only(ctx)
	require.NoError(t, err)
	assert.Len(t, parentWithChildren.Edges.Children, 1)
	assert.Equal(t, childTask.ID, parentWithChildren.Edges.Children[0].ID)

	t.Log("✅ Task hierarchy relationships working correctly")
}
