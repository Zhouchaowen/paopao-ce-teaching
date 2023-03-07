package conf

import (
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
