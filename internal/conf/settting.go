package conf

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

type ServerSettingS struct {
	RunMode      string
	HttpIp       string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type MySQLSettingS struct {
	UserName     string
	Password     string
	Host         string
	DBName       string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

func (s *MySQLSettingS) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		s.UserName,
		s.Password,
		s.Host,
		s.DBName,
		s.Charset,
		s.ParseTime,
	)
}

func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath(".")
	vp.AddConfigPath("configs/")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Setting{vp}, nil
}

func (s *Setting) Unmarshal(objects map[string]interface{}) error {
	for k, v := range objects {
		err := s.vp.UnmarshalKey(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
