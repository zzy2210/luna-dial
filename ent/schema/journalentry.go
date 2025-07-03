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

// JournalEntry holds the schema definition for the JournalEntry entity.
type JournalEntry struct {
	ent.Schema
}

// Fields of the JournalEntry.
func (JournalEntry) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id"),
		field.Text("content").
			NotEmpty(),
		field.String("time_reference").
			MaxLen(50).
			NotEmpty().
			Comment("时间引用，如'2025年7月'或'2025年W28周'"),
		field.Enum("time_scale").
			Values(
				string(types.TimeScaleDay),
				string(types.TimeScaleWeek),
				string(types.TimeScaleMonth),
				string(types.TimeScaleQuarter),
				string(types.TimeScaleYear),
			).
			Default(string(types.TimeScaleDay)),
		field.Enum("entry_type").
			Values(
				string(types.EntryTypePlanStart),
				string(types.EntryTypeReflection),
				string(types.EntryTypeSummary),
			).
			Default(string(types.EntryTypeReflection)),
		field.UUID("user_id", uuid.UUID{}),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the JournalEntry.
func (JournalEntry) Edges() []ent.Edge {
	return []ent.Edge{
		// 日志条目属于一个用户
		edge.From("user", User.Type).
			Ref("journal_entries").
			Field("user_id").
			Required().
			Unique(),
		// 日志条目可以关联多个任务
		edge.To("tasks", Task.Type),
	}
}

// Indexes of the JournalEntry.
func (JournalEntry) Indexes() []ent.Index {
	return []ent.Index{
		// 用户ID索引
		index.Fields("user_id"),
		// 时间尺度索引
		index.Fields("time_scale"),
		// 条目类型索引
		index.Fields("entry_type"),
		// 时间引用索引
		index.Fields("time_reference"),
		// 复合索引：用户+时间尺度
		index.Fields("user_id", "time_scale"),
		// 复合索引：用户+条目类型
		index.Fields("user_id", "entry_type"),
		// 创建时间索引
		index.Fields("created_at"),
	}
}
