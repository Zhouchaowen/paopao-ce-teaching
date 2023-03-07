package conf

import "log"

var (
	ServerSetting *ServerSettingS

	MysqlSetting *MySQLSettingS
)

func setupSetting() error {
	setting, err := NewSetting()
	if err != nil {
		return err
	}

	objects := map[string]interface{}{
		"Server": &ServerSetting,
		"MySQL":  &MysqlSetting,
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
