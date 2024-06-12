package middleware

import (
	"DisHub/common"
	"DisHub/common/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckMagicToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if CheckMagicTokenFromHeader(c) {
			// 通过
			c.Next() // 继续执行后续的处理函数
			return
		}
		// 不通过返回未授权
		responseAbortUnauthorized(c, "magic token not match")
	}
}

func CheckJWTToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(common.JWT_TOKEN)
		// 将令牌重新添加到请求头中
		//c.Request.Header.Set(common.JWT_TOKEN, token)
		if CheckMagicTokenFromHeader(c) {
			// 通过
			c.Next() // 继续执行后续的处理函数
			return
		}
		if token == "" {
			responseAbortUnauthorized(c, "jwt token empty")
			return
		}
		user, err := jwt.ParseUserToken(token)
		if err != nil {
			responseAbortUnauthorized(c, "jwt token expired")
			return
		}
		if user == nil {
			responseAbortUnauthorized(c, "jwt token invalided")
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "*")
		// 处理 OPTIONS 请求，直接返回 200 OK
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func CheckMagicTokenFromHeader(c *gin.Context) bool {
	token := c.Request.Header.Get(common.RPC_TOKEN_KEY)
	return token == common.RPC_TOKEN_VALVE
}

func responseAbortUnauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code": http.StatusUnauthorized,
		"msg":  msg,
	})
}
