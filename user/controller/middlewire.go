package controller

import (
	"net/http"
	"task-board/user/model"

	"task-board/jwt"

	"github.com/gin-gonic/gin"
)

func (u *UserController) MiddleFunc(c *gin.Context) gin.HandlerFunc {
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
					c.JSON(http.StatusOK, gin.H{"status": http.StatusBadGateway})
					c.Abort()
					return
				}
				if password != claimsmap["password"] {
					c.JSON(http.StatusBadRequest, gin.H{
						"status": http.StatusBadRequest,
						"msg":    "Wrong Password"})
					c.Abort()
					return
				}
			}
		}
	}
}
