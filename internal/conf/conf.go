package conf

import "log"

var (
	ServerSetting *ServerSettingS
)

func setupSetting() error {
	setting, err := NewSetting()
	if err != nil {
		return err
	}

	objects := map[string]interface{}{
		"Server": &ServerSetting,
	}
	if err = setting.Unmarshal(objects); err != nil {
		return err
	}

	return nil
}

func Initialize() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
}
