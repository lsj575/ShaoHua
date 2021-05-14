package initialize

import (
	"api/user-web/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("./user-web/config-dev.yaml") // 指定配置文件的文件名称
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			zap.S().Panic("config file not found", err.Error())
		} else {
			zap.S().Panic("config cannot load", err.Error())
		}
	}
	if err := v.Unmarshal(global.ServiceConfig); err != nil {
		zap.S().Panic("config unmarshal to struct failed", err.Error())
	}
	zap.S().Info("config loaded")
	zap.S().Infof("config info: %v", global.ServiceConfig)

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Info("config file chang ed: " + in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServiceConfig)
		zap.S().Info("config reloaded")
		zap.S().Infof("config info: %v", global.ServiceConfig)
	})
}
