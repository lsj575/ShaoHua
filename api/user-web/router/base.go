package router

import (
	"api/user-web/api"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	baseGroup := Router.Group("base")
	{
		baseGroup.POST("send_email_code", api.SendEmailVerificationCode)
	}
}
