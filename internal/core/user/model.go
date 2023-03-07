package user

import "gorm.io/plugin/soft_delete"

const (
	UserStatusNormal int = iota + 1
	UserStatusClosed
)

type User struct {
	ID         int64                 `gorm:"primary_key" json:"id"`
	Nickname   string                `json:"nickname"`
	Username   string                `json:"username"`
	Phone      string                `json:"phone"`
	Password   string                `json:"password"`
	Salt       string                `json:"salt"`
	Status     int                   `json:"status"`
	Avatar     string                `json:"avatar"`
	Balance    int64                 `json:"balance"`
	IsAdmin    bool                  `json:"is_admin"`
	CreatedOn  int64                 `json:"created_on"`
	ModifiedOn int64                 `json:"modified_on"`
	DeletedOn  int64                 `json:"deleted_on"`
	IsDel      soft_delete.DeletedAt `gorm:"softDelete:flag" json:"is_del"`
}
