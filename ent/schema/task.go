package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"

	"okr-web/internal/types"
)

// Task holds the schema definition for the Task entity.
type Task struct {
	ent.Schema
}

// Fields of the Task.
func (Task) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id"),
		field.String("title").
			MaxLen(255).
			NotEmpty(),
		field.Text("description").
			Optional(),
		field.Enum("type").
			Values(
				string(types.TaskTypeYear),
				string(types.TaskTypeQuarter),
				string(types.TaskTypeMonth),
				string(types.TaskTypeWeek),
				string(types.TaskTypeDay),
			).
			Default(string(types.TaskTypeDay)),
		field.Time("start_date").
			Default(time.Now),
		field.Time("end_date").
			Default(time.Now),
		field.Enum("status").
			Values(
				string(types.TaskStatusPending),
				string(types.TaskStatusInProgress),
				string(types.TaskStatusCompleted),
			).
			Default(string(types.TaskStatusPending)),
		field.Int("score").
			Default(0).
			Min(0).
			Max(10),
		field.UUID("parent_id", uuid.UUID{}).
			Optional().
			Nillable(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("tags").
			Optional().
			MaxLen(255),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Task.
func (Task) Edges() []ent.Edge {
	return []ent.Edge{
		// 任务属于一个用户
		edge.From("user", User.Type).
			Ref("tasks").
			Field("user_id").
			Required().
			Unique(),
		// 任务可以有父任务（自关联）
		edge.To("children", Task.Type).
			From("parent").
			Field("parent_id").
			Unique(),
		// 任务可以关联多个日志条目
		edge.From("journal_entries", JournalEntry.Type).
			Ref("tasks"),
	}
}

// Indexes of the Task.
func (Task) Indexes() []ent.Index {
	return []ent.Index{
		// 用户ID索引
		index.Fields("user_id"),
		// 父任务ID索引
		index.Fields("parent_id"),
		// 任务类型索引
		index.Fields("type"),
		// 开始时间索引
		index.Fields("start_date"),
		// 结束时间索引
		index.Fields("end_date"),
		// 状态索引
		index.Fields("status"),
		// 复合索引：用户+类型
		index.Fields("user_id", "type"),
		// 复合索引：用户+状态
		index.Fields("user_id", "status"),
	}
}
