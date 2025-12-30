package mongomodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`                                   // 标签名称
	Description string             `bson:"description,omitempty" json:"description,omitempty"` // 标签描述
	PostCount   int64              `bson:"post_count" json:"post_count"`                       // 使用该标签的帖子数
	IsOfficial  bool               `bson:"is_official" json:"is_official"`                     // 是否为官方标签
	CreateTime  uint               `bson:"create_time" json:"create_time"`                     // 创建时间
	UpdateTime  uint               `bson:"update_time" json:"update_time"`                     // 更新时间
	Status      int                `bson:"status" json:"status"`                               // 状态：1-正常，0-删除
}

func (t *Tag) CollectionName() string {
	return "tag"
}

const (
	TagFieldID         = "_id"
	TagFieldName       = "name"
	TagFieldPostCount  = "post_count"
	TagFieldIsOfficial = "is_official"
	TagFieldStatus     = "status"
	TagFieldCreateTime = "create_time"
	TagFieldUpdateTime = "update_time"
)

const (
	TagStatusNormal  = 1
	TagStatusDeleted = 0
)
