package main

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"srvs/user_srv/global"
	"srvs/user_srv/initialize"
	"srvs/user_srv/proto"
	"srvs/user_srv/services"
	"syscall"
)

func gracefulExit(server *grpc.Server) {
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	server.GracefulStop()
	log.Println("Server exiting")

}

func main() {
	// 1. initialize logger
	initialize.InitLogger()

	// 2. initialize config
	initialize.InitConfig()

	// 3.initialize mysql
	initialize.InitMysqlConnect(&global.ServiceConfig.DatabaseConfig.MysqlConfig)

	// 4. start server
	g := grpc.NewServer()
	proto.RegisterUserServer(g, services.UserService)
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		zap.S().Panic("failed to listen tcp:" + err.Error())
	}
	go func() {
		err = g.Serve(listener)
		if err != nil {
			zap.S().Panic("failed to start grpc:" + err.Error())
		}
	}()

	gracefulExit(g)
}
