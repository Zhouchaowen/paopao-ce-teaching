package post

import (
	"gorm.io/gorm"
	"time"
)

func CreatePost(db *gorm.DB, p *Post) (*Post, error) {
	p.LatestRepliedOn = time.Now().Unix()
	err := db.Create(&p).Error

	return p, err
}

func GetPostById(db *gorm.DB, id int64) (*Post, error) {
	var post Post
	if id > 0 {
		db = db.Where("id = ? AND is_del = ?", id, 0)
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	err := db.First(&post).Error
	if err != nil {
		return &post, err
	}

	return &post, nil
}

func DeletePostByPostId(db *gorm.DB, postId int64) error {
	return db.Model(&Post{}).Where("id = ?", postId).Updates(map[string]interface{}{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}
