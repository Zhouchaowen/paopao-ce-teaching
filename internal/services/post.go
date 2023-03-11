package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"paopao-ce-teaching/internal/conf"
	"paopao-ce-teaching/internal/cores/post"
	"paopao-ce-teaching/pkg/util"
	"strings"
)

type PostContentItem struct {
	Content string        `json:"content"  binding:"required"`
	Type    post.ContentT `json:"type"  binding:"required"`
	Sort    int64         `json:"sort"  binding:"required"`
}

// Check 检查PostContentItem属性
func (p *PostContentItem) Check() error {
	// 检查链接是否合法
	if p.Type == post.CONTENT_TYPE_LINK {
		if strings.Index(p.Content, "http://") != 0 && strings.Index(p.Content, "https://") != 0 {
			return fmt.Errorf("链接不合法")
		}
	}

	return nil
}

type PostCreationReq struct {
	Contents        []*PostContentItem `json:"contents" binding:"required"`
	Tags            []string           `json:"tags" binding:"required"`
	Users           []string           `json:"users" binding:"required"`
	AttachmentPrice int64              `json:"attachment_price"`
	Visibility      post.VisibleT      `json:"visibility"`
}

// CreatePost 创建文章
func CreatePost(c *gin.Context, userID int64, param PostCreationReq) (_ *post.Formatted, err error) {
	ip := c.ClientIP()
	tags := tagsFrom(param.Tags)
	postTmp := &post.Post{
		UserID:          userID,
		Tags:            strings.Join(tags, ","),
		IP:              ip,
		IPLoc:           util.GetIPLoc(ip),
		AttachmentPrice: param.AttachmentPrice,
		Visibility:      param.Visibility,
	}

	postTmp, err = post.CreatePost(conf.DB, postTmp)
	if err != nil {
		return nil, err
	}

	for _, item := range param.Contents {
		if err := item.Check(); err != nil {
			// 属性非法
			log.Infof("contents check err: %v", err)
			continue
		}

		if item.Type == post.CONTENT_TYPE_ATTACHMENT && param.AttachmentPrice > 0 {
			item.Type = post.CONTENT_TYPE_CHARGE_ATTACHMENT
		}

		postContent := &post.PostContent{
			PostID:  postTmp.ID,
			UserID:  userID,
			Content: item.Content,
			Type:    item.Type,
			Sort:    item.Sort,
		}
		if _, err = post.CreatePostContent(conf.DB, postContent); err != nil {
			return nil, err
		}
	}

	return postTmp.Format(), nil
}

func tagsFrom(originTags []string) []string {
	tags := make([]string, 0, len(originTags))
	for _, tag := range originTags {
		// TODO: 优化tag有效性检测
		if tag = strings.TrimSpace(tag); len(tag) > 0 {
			tags = append(tags, tag)
		}
	}
	return tags
}
