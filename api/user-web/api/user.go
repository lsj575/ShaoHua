package api

import (
	"api/user-web/api/splice"
	"api/user-web/global"
	"api/user-web/proto"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

func GetUserById(ctx *gin.Context) {
	idStr := ctx.DefaultQuery("id", "0")
	if idStr == "0" {
		ctx.JSON(http.StatusBadRequest, global.JsonError(global.ClientApiParamsError, "No legal parameters are passed in."))
		return
	}
	idInt, _ := strconv.Atoi(idStr)
	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		global.ServiceConfig.UserSrvConfig.Host, global.ServiceConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserById] dial error")
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	// 生成client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	rsp, err := userSrvClient.GetUserById(context.Background(), &proto.IDRequest{Id: uint64(idInt)})
	if err != nil {
		zap.S().Errorw("[GetUserById] get error")
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, global.JsonSuccess(fmt.Sprintf("Successfully obtain user %d information", idInt),
		splice.UserInfo(rsp)))
}

func GetUserByEmail(ctx *gin.Context) {
	email := ctx.DefaultQuery("email", "")
	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		global.ServiceConfig.UserSrvConfig.Host, global.ServiceConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserByEmail] dial error")
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	// 生成client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	rsp, err := userSrvClient.GetUserByEmail(context.Background(), &proto.EmailRequest{Email: email})
	if err != nil {
		zap.S().Errorw("[GetUserByEmail] get error")
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, global.JsonSuccess(fmt.Sprintf("Successfully obtain user %s information", email),
		splice.UserInfo(rsp)))
}

func GetCurrentUser(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	userIdInt, _ := userId.(uint64)
	if userIdInt <= 0 {
		ctx.JSON(http.StatusBadRequest, global.JsonError(global.ClientApiParamsError, "No legal parameters are passed in."))
		return
	}

	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		global.ServiceConfig.UserSrvConfig.Host, global.ServiceConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserById] dial error")
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	// 生成client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	rsp, err := userSrvClient.GetUserById(context.Background(), &proto.IDRequest{Id: uint64(userIdInt)})
	if err != nil {
		zap.S().Errorf("[GetCurrentUser] get error")
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, global.JsonSuccess(fmt.Sprintf("Successfully obtain user %d information", userIdInt),
		splice.UserInfo(rsp)))
}

func CheckUserExistByEmail(ctx *gin.Context) {
	email := ctx.DefaultQuery("email", "")
	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		global.ServiceConfig.UserSrvConfig.Host, global.ServiceConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[CheckUserExistByEmail] dial error")
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	// 生成client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	_, err = userSrvClient.GetUserByEmail(context.Background(), &proto.EmailRequest{Email: email})
	if err != nil {
		ctx.JSON(http.StatusOK, global.JsonSuccess(fmt.Sprintf("用户不存在"), false))
		return
	}
	ctx.JSON(http.StatusOK, global.JsonSuccess(fmt.Sprintf("用户已存在"), true))
}
