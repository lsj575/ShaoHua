package middlewares

import (
	"api/user-web/global"
	"api/user-web/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.Role != 1 {
			ctx.JSON(http.StatusForbidden, global.JsonError(global.ClientErrorRequest, "无访问权限"))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
