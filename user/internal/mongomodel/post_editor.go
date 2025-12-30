package mongomodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostEditor 文章编辑器模型
type PostEditor struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     uint64             `bson:"user_id" json:"user_id"` // 用户ID
	Title      string             `bson:"title" json:"title"`     // 标题
	Content    string             `bson:"content" json:"content"` // 内容
	Tags       []string           `bson:"tags" json:"tags"`
	Images     []string           `bson:"images" json:"images"`           // 图片URL列表
	UpdateTime uint               `bson:"update_time" json:"update_time"` // 更新时间
}

func (p *PostEditor) CollectionName() string {
	return "post_editor"
}

const (
	PostEditorFieldID         = "_id"
	PostEditorFieldUserID     = "user_id"
	PostEditorFieldTitle      = "title"
	PostEditorFieldContent    = "content"
	PostEditorFieldTags       = "tags"
	PostEditorFieldImageUrls  = "images"
	PostEditorFieldUpdateTime = "update_time"
)
