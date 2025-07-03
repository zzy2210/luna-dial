package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			StorageKey("id"),
		field.String("username").
			MaxLen(50).
			Unique().
			NotEmpty(),
		field.String("password").
			MaxLen(100).
			NotEmpty().
			Sensitive(), // 敏感字段，不会在序列化时显示
		field.String("email").
			MaxLen(100).
			Optional().
			Unique(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// 用户拥有多个任务
		edge.To("tasks", Task.Type),
		// 用户拥有多个日志条目
		edge.To("journal_entries", JournalEntry.Type),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// 用户名索引
		index.Fields("username").
			Unique(),
		// 邮箱索引
		index.Fields("email").
			Unique(),
		// 创建时间索引
		index.Fields("created_at"),
	}
}
