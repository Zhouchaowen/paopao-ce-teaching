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

func GetPostsByConditions(db *gorm.DB, conditions *map[string]interface{}, offset, limit int) ([]*Post, error) {
	var posts []*Post
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

	if err = db.Where("is_del = ?", 0).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

func CountByConditions(db *gorm.DB, conditions *map[string]interface{}) (int64, error) {
	var count int64

	for k, v := range *conditions {
		if k != "ORDER" {
			db = db.Where(k, v)
		}
	}
	if err := db.Model(&Post{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}
