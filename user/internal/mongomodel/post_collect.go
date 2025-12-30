package mongomodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostCollect struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID     primitive.ObjectID `bson:"post_id" json:"post_id"`         // 帖子ID
	UserID     uint64             `bson:"user_id" json:"user_id"`         // 用户ID
	CreateTime uint               `bson:"create_time" json:"create_time"` // 收藏时间
}

func (p *PostCollect) CollectionName() string {
	return "post_collect"
}

const (
	PostCollectFieldID         = "_id"
	PostCollectFieldPostID     = "post_id"
	PostCollectFieldUserID     = "user_id"
	PostCollectFieldCreateTime = "create_time"
)
