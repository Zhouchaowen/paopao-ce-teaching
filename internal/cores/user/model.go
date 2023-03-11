package user

import "gorm.io/plugin/soft_delete"

const (
	UserStatusNormal int = iota + 1
	UserStatusClosed
)

type User struct {
	ID         int64                 `gorm:"primary_key" json:"id"`         // 用户ID
	Nickname   string                `json:"nickname"`                      // 昵称
	Username   string                `json:"username"`                      // 用户名
	Phone      string                `json:"phone"`                         // 手机号
	Password   string                `json:"password"`                      // MD5密码
	Salt       string                `json:"salt"`                          // 盐值
	Status     int                   `json:"status"`                        // 状态，1正常，2停用
	Avatar     string                `json:"avatar"`                        // 用户头像
	Balance    int64                 `json:"balance"`                       // 用户余额（分）
	IsAdmin    bool                  `json:"is_admin"`                      // 是否管理员
	CreatedOn  int64                 `json:"created_on"`                    // 创建时间
	ModifiedOn int64                 `json:"modified_on"`                   // 修改时间
	DeletedOn  int64                 `json:"deleted_on"`                    // 删除时间
	IsDel      soft_delete.DeletedAt `gorm:"softDelete:flag" json:"is_del"` // 是否删除 0 为未删除、1 为已删除
}

type Formatted struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Status   int    `json:"status"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`
}
