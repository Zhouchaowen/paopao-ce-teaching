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
