package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Carts holds the schema definition for the Carts entity.
type Carts struct {
	ent.Schema
}

// Fields of the Carts.
func (Carts) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("cart_id").Immutable().Unique().Comment("购物车项ID").StructTag(`json:"cart_id"`),
		field.Int64("user_id").Comment("用户ID").StructTag(`json:"user_id"`),
		field.Int64("product_id").Comment("商品ID").StructTag(`json:"product_id"`),
		field.Int("quantity").Default(1).Comment("选购数量").StructTag(`json:"quantity"`),
		field.Int("status").Optional().Default(1).Comment("状态: 1-正常, 2-已删除").StructTag(`json:"status"`),
		field.Int64("company_id").Comment("公司ID").StructTag(`json:"company_id"`),
		field.String("company_name").Comment("公司名称").StructTag(`json:"company_name"`).MaxLen(255),
		field.Time("created_at").Optional().Comment("创建时间").StructTag(`json:"created_at"`),
		field.Time("updated_at").Optional().Comment("最后更新时间").StructTag(`json:"updated_at"`),
	}
}

// Edges of the Carts.
func (Carts) Edges() []ent.Edge {
	return []ent.Edge{}
}

func (Carts) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "carts"},
	}
}
