package post

import (
	"gorm.io/plugin/soft_delete"
	"paopao-ce-teaching/internal/cores/user"
	"strings"
)

// VisibleT 可访问类型，0公开，1私密，2好友
type VisibleT uint8

const (
	VisitPublic VisibleT = iota
	VisitPrivate
	VisitFriend
	VisitInvalid
)

type Post struct {
	ID              int64                 `gorm:"primary_key" json:"id"`         // 主题ID
	UserID          int64                 `json:"user_id"`                       // 用户ID
	CommentCount    int64                 `json:"comment_count"`                 // 评论数
	CollectionCount int64                 `json:"collection_count"`              // 收藏数
	UpvoteCount     int64                 `json:"upvote_count"`                  // 点赞数
	Visibility      VisibleT              `json:"visibility"`                    // 可见性 0公开 1私密 2好友可见
	IsTop           int                   `json:"is_top"`                        // 是否置顶
	IsEssence       int                   `json:"is_essence"`                    // 是否精华
	IsLock          int                   `json:"is_lock"`                       // 是否锁定
	LatestRepliedOn int64                 `json:"latest_replied_on"`             // 最新回复时间
	Tags            string                `json:"tags"`                          // 标签
	AttachmentPrice int64                 `json:"attachment_price"`              // 附件价格(分)
	IP              string                `json:"ip"`                            // IP地址
	IPLoc           string                `json:"ip_loc"`                        // IP城市地址
	CreatedOn       int64                 `json:"created_on"`                    // 创建时间
	ModifiedOn      int64                 `json:"modified_on"`                   // 修改时间
	DeletedOn       int64                 `json:"deleted_on"`                    // 删除时间
	IsDel           soft_delete.DeletedAt `gorm:"softDelete:flag" json:"is_del"` // 是否删除 0 为未删除、1 为已删除
}

func (p *Post) Format() *Formatted {
	if p.ID > 0 {
		tagsMap := map[string]int8{}
		for _, tag := range strings.Split(p.Tags, ",") {
			tagsMap[tag] = 1
		}
		return &Formatted{
			ID:              p.ID,
			UserID:          p.UserID,
			User:            &user.Formatted{},
			Contents:        []*ContentFormatted{},
			CommentCount:    p.CommentCount,
			CollectionCount: p.CollectionCount,
			UpvoteCount:     p.UpvoteCount,
			Visibility:      p.Visibility,
			IsTop:           p.IsTop,
			IsEssence:       p.IsEssence,
			IsLock:          p.IsLock,
			LatestRepliedOn: p.LatestRepliedOn,
			CreatedOn:       p.CreatedOn,
			ModifiedOn:      p.ModifiedOn,
			AttachmentPrice: p.AttachmentPrice,
			Tags:            tagsMap,
			IPLoc:           p.IPLoc,
		}
	}

	return nil
}

type Formatted struct {
	ID              int64               `json:"id"`
	UserID          int64               `json:"user_id"`
	User            *user.Formatted     `json:"user"`
	Contents        []*ContentFormatted `json:"contents"`
	CommentCount    int64               `json:"comment_count"`
	CollectionCount int64               `json:"collection_count"`
	UpvoteCount     int64               `json:"upvote_count"`
	Visibility      VisibleT            `json:"visibility"`
	IsTop           int                 `json:"is_top"`
	IsEssence       int                 `json:"is_essence"`
	IsLock          int                 `json:"is_lock"`
	LatestRepliedOn int64               `json:"latest_replied_on"`
	CreatedOn       int64               `json:"created_on"`
	ModifiedOn      int64               `json:"modified_on"`
	Tags            map[string]int8     `json:"tags"`
	AttachmentPrice int64               `json:"attachment_price"`
	IPLoc           string              `json:"ip_loc"`
}

type PostContent struct {
	ID         int64                 `gorm:"primary_key" json:"id"`         // 内容ID
	PostID     int64                 `json:"post_id"`                       // POST ID
	UserID     int64                 `json:"user_id"`                       // 用户ID
	Content    string                `json:"content"`                       // 内容
	Type       ContentT              `json:"type"`                          // 类型，1标题，2文字段落，3图片地址，4视频地址，5语音地址，6链接地址，7附件资源，8收费资源
	Sort       int64                 `json:"sort"`                          // 排序，越小越靠前
	CreatedOn  int64                 `json:"created_on"`                    // 创建时间
	ModifiedOn int64                 `json:"modified_on"`                   // 修改时间
	DeletedOn  int64                 `json:"deleted_on"`                    // 删除时间
	IsDel      soft_delete.DeletedAt `gorm:"softDelete:flag" json:"is_del"` // 是否删除 0 为未删除、1 为已删除
}

func (p *PostContent) Format() *ContentFormatted {
	if p.ID == 0 {
		return nil
	}
	return &ContentFormatted{
		ID:      p.ID,
		PostID:  p.PostID,
		Content: p.Content,
		Type:    p.Type,
		Sort:    p.Sort,
	}
}

// 类型，1标题，2文字段落，3图片地址，4视频地址，5语音地址，6链接地址，7附件资源

type ContentT int

const (
	CONTENT_TYPE_TITLE ContentT = iota + 1
	CONTENT_TYPE_TEXT
	CONTENT_TYPE_IMAGE
	CONTENT_TYPE_VIDEO
	CONTENT_TYPE_AUDIO
	CONTENT_TYPE_LINK
	CONTENT_TYPE_ATTACHMENT
	CONTENT_TYPE_CHARGE_ATTACHMENT
)

var (
	mediaContentType = []ContentT{
		CONTENT_TYPE_IMAGE,
		CONTENT_TYPE_VIDEO,
		CONTENT_TYPE_AUDIO,
		CONTENT_TYPE_ATTACHMENT,
		CONTENT_TYPE_CHARGE_ATTACHMENT,
	}
)

type ContentFormatted struct {
	ID      int64    `json:"id"`
	PostID  int64    `json:"post_id"`
	Content string   `json:"content"`
	Type    ContentT `json:"type"`
	Sort    int64    `json:"sort"`
}
