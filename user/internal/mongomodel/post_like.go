package mongomodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PostLike 用户对帖子的点赞记录
type PostLike struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID     primitive.ObjectID `bson:"post_id" json:"post_id"` // 帖子ID（关联post集合的_id）
	UserID     uint64             `bson:"user_id" json:"user_id"` // 点赞用户ID
	Status     int8               `bson:"status" json:"status"`
	CreateTime uint               `bson:"create_time" json:"create_time"` // 点赞时间
}

func (pl *PostLike) CollectionName() string {
	return "post_like"
}

const (
	PostLikeStatusNormal  = 1
	PostLikeStatusDeleted = 0
)

const (
	PostLikeFieldID         = "_id"
	PostLikeFieldPostID     = "post_id"
	PostLikeFieldUserID     = "user_id"
	PostLikeFieldStatus     = "status"
	PostLikeFieldCreateTime = "create_time"
)
