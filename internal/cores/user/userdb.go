package user

import (
	"gorm.io/gorm"
)

func GetUserByID(db *gorm.DB, id int64) (*User, error) {
	var user User
	db = db.Where("id= ? AND is_del = ?", id, 0)

	err := db.First(&user).Error
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	db = db.Where("username = ? AND is_del = ?", username, 0)

	err := db.First(&user).Error
	if err != nil {
		return &user, err
	}

	return &user, nil
}

func Create(db *gorm.DB, u *User) (*User, error) {
	err := db.Create(&u).Error

	return u, err
}
