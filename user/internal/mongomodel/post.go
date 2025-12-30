package mongomodel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post 帖子内容模型
type Post struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        uint64             `bson:"user_id" json:"user_id"`                                   // 发布者ID
	SchoolID      uint               `bson:"school_id" json:"school_id"`                               // 学校ID（关联学校表）
	Title         string             `bson:"title" json:"title"`                                       // 帖子标题
	Content       string             `bson:"content" json:"content"`                                   // 文本内容
	Images        []string           `bson:"images" json:"images"`                                     // 图片URL列表
	Tags          []string           `bson:"tags" json:"tags,omitempty"`                               // 标签（如："校招","技术面"）
	Anonymous     bool               `bson:"anonymous" json:"anonymous"`                               // 是否匿名发布
	AnonymousName string             `bson:"anonymous_name,omitempty" json:"anonymous_name,omitempty"` // 匿名称呼（如："匿名用户"）
	CreateTime    uint               `bson:"create_time" json:"create_time"`                           // 发布时间
	UpdateTime    uint               `bson:"update_time" json:"update_time"`                           // 更新时间
	Status        int                `bson:"status" json:"status"`                                     // 状态：1-正常，0-删除，2-审核中,3-审核失败,4-仅自己可见
	LikeCount     int64              `bson:"like_count" json:"like_count"`                             // 点赞数
	CommentCount  int64              `bson:"comment_count" json:"comment_count"`                       // 评论数
	CollectCount  int64              `bson:"collect_count" json:"collect_count"`                       // 收藏数
}

func (p *Post) CollectionName() string {
	return "post"
}

const (
	PostFieldID            = "_id"
	PostFieldUserID        = "user_id"
	PostFieldSchoolID      = "school_id"
	PostFieldContent       = "content"
	PostFieldImages        = "images"
	PostFieldTags          = "tags"
	PostFieldAnonymous     = "anonymous"
	PostFieldAnonymousName = "anonymous_name"
	PostFieldCreateTime    = "create_time"
	PostFieldUpdateTime    = "update_time"
	PostFieldStatus        = "status"
	PostFieldLikeCount     = "like_count"
	PostFieldCommentCount  = "comment_count"
	PostFieldCollectCount  = "collect_count"
)

const (
	PostStatusNormal  = 1
	PostStatusDeleted = 0
	PostStatusAudit   = 2
	PostStatusFail    = 3
	PostStatusPrivate = 4
)

const (
	PostTypeGeneral = 0 // 通用推荐
	PostTypeCampus  = 1 // 校园帖子（本校内容）
	PostTypeVote    = 2 // 投票帖子（包含"投票"标签）
	PostTypeRating  = 3 // 评分帖子（包含"评分"标签）
)
