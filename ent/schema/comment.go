package schema

import (
	"time"

	"cp-website/util"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Comment struct {
	ent.Schema
}

func (Comment) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").DefaultFunc(util.GenID).Immutable(),
		field.String("content").NotEmpty(),
		field.Time("created_at").Default(time.Now).Immutable(), // 自动设置当前时间
	}
}

func (Comment) Edges() []ent.Edge {
	return []ent.Edge{
		// 这条评论是谁发的
		edge.From("author", User.Type).Ref("comments").Unique().Required(),
		// 这条评论属于哪个 CP
		edge.From("cp", CP.Type).Ref("comments").Unique().Required(),

		// 🌟 【核心：自引用关系】用于嵌套评论
		// 父评论 (Parent)
		edge.To("children", Comment.Type).From("parent").Unique(),
	}
}
