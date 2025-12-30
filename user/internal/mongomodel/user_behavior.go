// 在 internal/mongomodel/ 目录下创建 user_behavior.go 文件
package mongomodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserBehavior 用户行为模型
type UserBehavior struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID     uint64             `bson:"user_id" json:"user_id"`         // 用户ID
	PostID     primitive.ObjectID `bson:"post_id" json:"post_id"`         // 帖子ID
	ActionType string             `bson:"action_type" json:"action_type"` // 行为类型: like, comment, collect, view
	Timestamp  uint               `bson:"timestamp" json:"timestamp"`     // 时间戳
}

func (u *UserBehavior) CollectionName() string {
	return "user_behavior"
}

const (
	UserBehaviorFieldID         = "_id"
	UserBehaviorFieldUserID     = "user_id"
	UserBehaviorFieldPostID     = "post_id"
	UserBehaviorFieldActionType = "action_type"
	UserBehaviorFieldTimestamp  = "timestamp"
)

const (
	ActionTypeLike    = "like"
	ActionTypeComment = "comment"
	ActionTypeCollect = "collect"
	ActionTypeView    = "view"
)
