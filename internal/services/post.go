package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"gorm.io/gorm"
	"paopao-ce-teaching/internal/conf"
	"paopao-ce-teaching/internal/cores/post"
	"paopao-ce-teaching/internal/cores/user"
	"paopao-ce-teaching/pkg/errors"
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

func GetPost(id int64) (*post.Formatted, error) {
	postTmp, err := post.GetPostById(conf.DB, id)
	if err != nil {
		return nil, err
	}

	postContents, err := post.GetPostContentsByConditions(
		conf.DB,
		&map[string]interface{}{
			"post_id = ?": postTmp.ID,
			"ORDER":       "sort ASC",
		}, 0, 0,
	)
	if err != nil {
		return nil, err
	}

	users, err := user.GetUsersByConditions(
		conf.DB,
		&map[string]interface{}{
			"id = ?": postTmp.UserID,
		}, 0, 0,
	)
	if err != nil {
		return nil, err
	}

	// 数据整合
	postFormatted := postTmp.Format()
	for _, user := range users {
		postFormatted.User = user.Format()
	}
	for _, content := range postContents {
		if content.PostID == postTmp.ID {
			postFormatted.Contents = append(postFormatted.Contents, content.Format())
		}
	}
	return postFormatted, nil
}

type PostDelReq struct {
	ID int64 `json:"id" binding:"required"`
}

func DeletePost(userId int64, id int64) *errors.Error {
	postTmp, err := post.GetPostById(conf.DB, id)
	if err != nil {
		return errors.GetPostFailed
	}

	userTmp, err := user.GetUserByID(conf.DB, userId)
	if err != nil {
		return errors.NoPermission
	}
	if postTmp.UserID != userId && !userTmp.IsAdmin {
		return errors.NoPermission
	}

	err = conf.DB.Transaction(
		func(tx *gorm.DB) error {

			// 删推文
			if err := post.DeletePostByPostId(tx, postTmp.ID); err != nil {
				return err
			}

			// 删内容
			if err := post.DeletePostContentByPostId(tx, postTmp.ID); err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		log.Errorf("service.DeletePost delete post failed: %s", err)
		return errors.DeletePostFailed
	}

	return nil
}

type PostListReq struct {
	Conditions *map[string]interface{}
	Offset     int
	Limit      int
}

func GetPostList(req *PostListReq) ([]*post.Formatted, error) {
	posts, err := post.GetPostsByConditions(conf.DB, req.Conditions, req.Offset, req.Limit)

	if err != nil {
		return nil, err
	}

	return MergePosts(posts)
}

// MergePosts post数据整合
func MergePosts(posts []*post.Post) ([]*post.Formatted, error) {
	postIds := make([]int64, 0, len(posts))
	userIds := make([]int64, 0, len(posts))
	for _, post := range posts {
		postIds = append(postIds, post.ID)
		userIds = append(userIds, post.UserID)
	}

	postContents, err := post.GetPostContentsByConditions(
		conf.DB,
		&map[string]interface{}{
			"post_id IN ?": postIds,
			"ORDER":        "sort ASC",
		}, 0, 0,
	)
	if err != nil {
		return nil, err
	}

	users, err := user.GetUsersByConditions(
		conf.DB,
		&map[string]interface{}{
			"id IN ?": userIds,
		}, 0, 0,
	)
	if err != nil {
		return nil, err
	}

	userMap := make(map[int64]*user.Formatted, len(users))
	for _, user := range users {
		userMap[user.ID] = user.Format()
	}

	contentMap := make(map[int64][]*post.ContentFormatted, len(postContents))
	for _, content := range postContents {
		contentMap[content.PostID] = append(contentMap[content.PostID], content.Format())
	}

	// 数据整合
	postsFormatted := make([]*post.Formatted, 0, len(posts))
	for _, post := range posts {
		postFormatted := post.Format()
		postFormatted.User = userMap[post.UserID]
		postFormatted.Contents = contentMap[post.ID]
		postsFormatted = append(postsFormatted, postFormatted)
	}
	return postsFormatted, nil
}

func GetPostCount(conditions *map[string]interface{}) (int64, error) {
	return post.CountByConditions(conf.DB, conditions)
}
