package user

import (
	"gorm.io/gorm"
)

func Create(db *gorm.DB, u *User) (*User, error) {
	err := db.Create(&u).Error

	return u, err
}

func UpdateUser(db *gorm.DB, u *User) error {
	return db.Model(&User{}).Where("id = ? AND is_del = ?", u.ID, 0).Save(u).Error
}

func Get(db *gorm.DB, u *User) (*User, error) {
	var user User
	if u.ID > 0 {
		db = db.Where("id= ? AND is_del = ?", u.ID, 0)
	} else if u.Phone != "" {
		db = db.Where("phone = ? AND is_del = ?", u.Phone, 0)
	} else {
		db = db.Where("username = ? AND is_del = ?", u.Username, 0)
	}

	err := db.First(&user).Error
	if err != nil {
		return &user, err
	}

	return &user, nil
}

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

func GetUsersByConditions(db *gorm.DB, conditions *map[string]interface{}, offset, limit int) ([]*User, error) {
	var users []*User
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

	if err = db.Where("is_del = ?", 0).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
