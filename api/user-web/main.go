package main

import (
	"api/user-web/global"
	"api/user-web/initialize"
	myvalidator "api/user-web/validator"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. initialize logger
	initialize.InitLogger()

	// 2. initialize router
	router := initialize.Routers()

	// 3. initialize config
	initialize.InitConfig()

	// 4. register validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("loginEmail", myvalidator.ValidateLoginEmail)
		_ = v.RegisterValidation("loginUsername", myvalidator.ValidateLoginUsername)
	}

	go func() {
		zap.S().Infof("启动服务器, 端口：%d", global.ServiceConfig.Port)
		if err := router.Run(fmt.Sprintf(":%d", global.ServiceConfig.Port)); err != nil {
			zap.S().Panic("启动失败", err.Error())
		}
	}()

	// graceful exit
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.S().Info("start to stop user-web api server")
}
