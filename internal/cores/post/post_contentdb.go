package post

import (
	"gorm.io/gorm"
	"time"
)

func CreatePostContent(db *gorm.DB, p *PostContent) (*PostContent, error) {
	err := db.Create(&p).Error

	return p, err
}

func GetPostContentsByConditions(db *gorm.DB, conditions *map[string]interface{}, offset, limit int) ([]*PostContent, error) {
	var contents []*PostContent
	var err error
	if offset >= 0 && limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}

	for k, v := range *conditions {
		if k == "ORDER" {
			db = db.Order(v)
		} else {
			db = db.Where(k, v)
		}
	}

	if err = db.Where("is_del = ?", 0).Find(&contents).Error; err != nil {
		return nil, err
	}

	return contents, nil
}

func DeletePostContentByPostId(db *gorm.DB, postId int64) error {
	return db.Model(&PostContent{}).Where("post_id = ?", postId).Updates(map[string]interface{}{
		"deleted_on": time.Now().Unix(),
		"is_del":     1,
	}).Error
}
