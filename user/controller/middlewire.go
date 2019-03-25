package controller

import (
	"TEST/user/model"
	"net/http"

	"TEST/jwt"

	"github.com/gin-gonic/gin"
)

func (u *UserController) Jwttest(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("token")
		if tokenString == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "no token",
			})
			c.Abort()
			return
		}
		if claims, ok := jwt.ParseToken(tokenString); ok {
			if claimsmap, ok := claims.(map[string]string); ok {
				password, err := model.InfoPasswordByName(u.db, u.tableName, claimsmap["name"])
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"status": err})
					c.Abort()
					return
				}
				if password == claimsmap["password"] {
					c.JSON(http.StatusOK, gin.H{
						"status": http.StatusOK,
						"msg":    "密码错误"})
				}
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "err"})
			c.Abort()
			return
		}
	}
}
