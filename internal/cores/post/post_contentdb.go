package post

import "gorm.io/gorm"

func CreatePostContent(db *gorm.DB, p *PostContent) (*PostContent, error) {
	err := db.Create(&p).Error

	return p, err
}
