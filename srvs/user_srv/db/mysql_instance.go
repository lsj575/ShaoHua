package db

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"srvs/user_srv/config"
	"srvs/user_srv/models"
	"strconv"
	"sync"
)

type MysqlConnectPool struct {
	db *gorm.DB
}

var (
	instance *MysqlConnectPool
	once     sync.Once
	errDb    error
)

func GetMysqlInstance() *MysqlConnectPool {
	once.Do(func() {
		instance = &MysqlConnectPool{}
	})
	return instance
}

func (m *MysqlConnectPool) InitDataPool(mysqlConfig *config.MysqlConfig) bool {
	dsn := mysqlConfig.Username + ":" + mysqlConfig.Password + "@tcp(" + mysqlConfig.Host + ":" +
		strconv.Itoa(mysqlConfig.Port) + ")/" + mysqlConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	m.db, errDb = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	_ = m.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&models.User{})
	if errDb != nil {
		zap.S().Panic(errDb)
		return false
	}
	return true
}

func (m *MysqlConnectPool) GetMysqlDB() *gorm.DB {
	return m.db
}
