package mongomodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID     primitive.ObjectID `bson:"post_id" json:"post_id"`                         // 关联的帖子ID
	UserID     uint64             `bson:"user_id" json:"user_id"`                         // 评论者ID
	Content    string             `bson:"content" json:"content"`                         // 评论内容
	ParentID   primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"` // 父评论ID（用于回复）
	Anonymous  bool               `bson:"anonymous" json:"anonymous"`                     // 是否匿名评论
	CreateTime uint               `bson:"create_time" json:"create_time"`                 // 创建时间
	Status     int8               `bson:"status" json:"status"`                           // 状态：1-正常，0-删除
}

const (
	CommentStatusNormal  = 1
	CommentStatusDeleted = 0
)

const (
	CommentFieldID         = "_id"
	CommentFieldPostID     = "post_id"
	CommentFieldUserID     = "user_id"
	CommentFieldContent    = "content"
	CommentFieldParentID   = "parent_id"
	CommentFieldAnonymous  = "anonymous"
	CommentFieldCreateTime = "create_time"
	CommentFieldStatus     = "status"
)
