package initialize

import (
	"srvs/user_srv/config"
	"srvs/user_srv/db"
)

func InitMysqlConnect(mysqlConfig *config.MysqlConfig) {
	db.GetMysqlInstance().InitDataPool(mysqlConfig)

}
