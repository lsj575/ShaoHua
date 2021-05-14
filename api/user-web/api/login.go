package api

import (
	"api/user-web/forms"
	"api/user-web/global"
	"api/user-web/middlewares"
	"api/user-web/models"
	"api/user-web/proto"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"time"
)

func HandleValidatorError(ctx *gin.Context, err error) {
	_, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusBadRequest, global.JsonError(global.ClientErrorRequest, err.Error()))
		return
	}
	ctx.JSON(http.StatusBadRequest, global.JsonError(global.ClientErrorRequest, "Invalid Argument "+err.Error()))
}

func PasswordLogin(ctx *gin.Context) {
	passwordLoginForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		global.ServiceConfig.UserSrvConfig.Host, global.ServiceConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[PasswordLogin] dial error")
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	// 生成client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	// login logic
	if res, err := userSrvClient.GetUserByEmail(context.Background(), &proto.EmailRequest{
		Email: passwordLoginForm.Email,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, global.JsonError(global.ClientGRPCParamsError, "用户不存在"))
			default:
				ctx.JSON(http.StatusInternalServerError, global.JsonError(global.ServerGRPCInternalError, "登录失败"))
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, global.JsonError(global.ServerGRPCUnknownError, "登录失败"))
		}
		return
	} else {
		if passRes, passErr := userSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.Password,
			EncryptedPassword: res.Password,
		}); passErr != nil {
			ctx.JSON(http.StatusInternalServerError, global.JsonError(global.ServerGRPCInternalError, "登录失败"))
		} else {
			if passRes.Success {
				// generate token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:       res.Id,
					Username: res.Username,
					Role:     res.Roles,
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               // 签名生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30天过期
						Issuer:    "token",
					},
				}
				token, tokenErr := j.CreateToken(claims)
				if tokenErr != nil {
					ctx.JSON(http.StatusInternalServerError, global.JsonError(global.ServerAPIInternalError, "生成token失败"))
				}

				ctx.JSON(http.StatusOK, global.JsonSuccess("登录成功", gin.H{
					"id":         res.Id,
					"username":   res.Username,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				}))
			} else {
				ctx.JSON(http.StatusBadRequest, global.JsonSuccess("密码错误", []int{}))
			}
		}
	}
}

/**
用户免密登录
*/
func NonPasswordLogin(ctx *gin.Context) {
	nonPasswordLoginForm := forms.NonPasswordLoginForm{}
	if err := ctx.ShouldBind(&nonPasswordLoginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// verification code check
	if !VerifyString(nonPasswordLoginForm.Email, nonPasswordLoginForm.Code) {
		ctx.JSON(http.StatusBadRequest, global.JsonError(global.ClientApiParamsError, "验证码错误"))
		return
	}

	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		global.ServiceConfig.UserSrvConfig.Host, global.ServiceConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[NonPasswordLogin] dial error")
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	// 生成client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	// login logic
	if res, err := userSrvClient.GetUserByEmail(context.Background(), &proto.EmailRequest{
		Email: nonPasswordLoginForm.Email,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, global.JsonError(global.ClientGRPCParamsError, "用户不存在"))
			default:
				ctx.JSON(http.StatusInternalServerError, global.JsonError(global.ServerGRPCInternalError, "登录失败"))
			}
		} else {
			ctx.JSON(http.StatusInternalServerError, global.JsonError(global.ServerGRPCUnknownError, "登录失败"))
		}
		return
	} else {
		// 因为之前验证过验证码，所以不需要重新验证
		// generate token
		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			ID:       res.Id,
			Username: res.Username,
			Role:     res.Roles,
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),               // 签名生效时间
				ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30天过期
				Issuer:    "token",
			},
		}
		token, tokenErr := j.CreateToken(claims)
		if tokenErr != nil {
			ctx.JSON(http.StatusInternalServerError, global.JsonError(global.ServerAPIInternalError, "生成token失败"))
		}

		ctx.JSON(http.StatusOK, global.JsonSuccess("登录成功", gin.H{
			"id":         res.Id,
			"username":   res.Username,
			"token":      token,
			"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
		}))
	}
}

/**
用户注册
*/
func Register(ctx *gin.Context) {
	registerForm := forms.RegisterForm{}
	if err := ctx.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// verification code check
	if !VerifyString(registerForm.Email, registerForm.Code) {
		ctx.JSON(http.StatusBadRequest, global.JsonError(global.ClientApiParamsError, "验证码错误"))
		return
	}

	// 拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		global.ServiceConfig.UserSrvConfig.Host, global.ServiceConfig.UserSrvConfig.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[Register] dial error")
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	// 生成client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	userInfo, err := userSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Username:   "小吾" + registerForm.Email[:6] + strconv.Itoa(time.Now().Minute()),
		Password:   "",
		RePassword: "",
		Email:      registerForm.Email,
	})

	if err != nil {
		zap.S().Errorf("[Register] get error", err.Error())
		global.HandleGrpcErrorToHttpError(err, ctx)
		return
	}

	// auto login after register
	// generate token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:       userInfo.Id,
		Username: userInfo.Username,
		Role:     userInfo.Roles,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, // 30天过期
			Issuer:    "token",
		},
	}
	token, tokenErr := j.CreateToken(claims)
	if tokenErr != nil {
		ctx.JSON(http.StatusInternalServerError, global.JsonError(global.ServerAPIInternalError, "生成token失败"))
	}

	ctx.JSON(http.StatusOK, global.JsonSuccess("注册成功", gin.H{
		"id":         userInfo.Id,
		"username":   userInfo.Username,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	}))
}
