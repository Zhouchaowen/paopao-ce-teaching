package services

import (
	"github.com/gofrs/uuid"
	"paopao-ce-teaching/internal/conf"
	"paopao-ce-teaching/internal/cores/user"
	"paopao-ce-teaching/pkg/errors"
	"paopao-ce-teaching/pkg/util"
	"regexp"
	"strings"
	"unicode/utf8"
)

type AuthRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type ChangePasswordReq struct {
	Password    string `json:"password" form:"password" binding:"required"`
	OldPassword string `json:"old_password" form:"old_password" binding:"required"`
}

type ChangeNicknameReq struct {
	Nickname string `json:"nickname" form:"nickname" binding:"required"`
}

// DoLogin 用户认证
func DoLogin(param *AuthRequest) (*user.User, error) {
	u, err := user.GetUserByUsername(conf.DB, param.Username)
	if err != nil {
		return nil, errors.UnauthorizedAuthNotExist
	}

	if u.ID > 0 {
		// 对比密码是否正确
		if ValidPassword(u.Password, param.Password, u.Salt) {

			if u.Status == user.Closed {
				return nil, errors.UserHasBeenBanned
			}

			return u, nil
		}

		return nil, errors.UnauthorizedAuthFailed
	}

	return nil, errors.UnauthorizedAuthNotExist
}

// ValidPassword 检查密码是否一致
func ValidPassword(dbPassword, password, salt string) bool {
	return strings.Compare(dbPassword, util.EncodeMD5(util.EncodeMD5(password)+salt)) == 0
}

// ValidUsername 验证用户
func ValidUsername(username string) error {
	// 检测用户是否合规
	if utf8.RuneCountInString(username) < 3 || utf8.RuneCountInString(username) > 12 {
		return errors.UsernameLengthLimit
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username) {
		return errors.UsernameCharLimit
	}

	// 重复检查
	user, _ := user.GetUserByUsername(conf.DB, username)

	if user.ID > 0 {
		return errors.UsernameHasExisted
	}

	return nil
}

// CheckPassword 密码检查
func CheckPassword(password string) error {
	// 检测用户是否合规
	if utf8.RuneCountInString(password) < 6 || utf8.RuneCountInString(password) > 16 {
		return errors.PasswordLengthLimit
	}

	return nil
}

// EncryptPasswordAndSalt 密码加密&生成salt
func EncryptPasswordAndSalt(password string) (string, string) {
	salt := uuid.Must(uuid.NewV4()).String()[:8]
	password = util.EncodeMD5(util.EncodeMD5(password) + salt)

	return password, salt
}

// Register 用户注册
func Register(username, password string) (*user.User, error) {
	password, salt := EncryptPasswordAndSalt(password)

	userTmp := &user.User{
		Nickname: username,
		Username: username,
		Password: password,
		Avatar:   "test.png",
		Salt:     salt,
		Status:   user.Normal,
	}

	userTmp, err := user.Create(conf.DB, userTmp)
	if err != nil {
		return nil, err
	}

	return userTmp, nil
}

// GetUserByUsername 通过用户名获取用户基本信息
func GetUserByUsername(username string) (*user.User, error) {
	u, err := user.GetUserByUsername(conf.DB, username)

	if err != nil {
		return nil, err
	}

	if u.ID > 0 {
		return u, nil
	}

	return nil, errors.NoExistUsername
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(u *user.User) *errors.Error {
	if err := user.UpdateUser(conf.DB, u); err != nil {
		return errors.ServerError
	}
	return nil
}
