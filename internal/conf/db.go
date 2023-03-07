package conf

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"sync"
)

var (
	DB   *gorm.DB
	once sync.Once
)

func MustGormDB() *gorm.DB {
	once.Do(func() {
		var err error
		if DB, err = newDBEngine(); err != nil {
			log.Fatalf("new gorm db failed: %s", err)
		}
	})

	return DB
}

func newDBEngine() (*gorm.DB, error) {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "p_",
			SingularTable: true,
		},
	}
	db, err := gorm.Open(mysql.Open(MysqlSetting.Dsn()), config)
	if err != nil {
		return nil, err
	}
	return db, nil
}
