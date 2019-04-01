package controller

import (
	"net/http"

	"task-board/jwt"

	"github.com/gin-gonic/gin"
)

func (u *UserController) MiddleFunc(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("token")
		if tokenString == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "no token",
			})
			c.Abort()
			return
		}
		claims, ok := jwt.ParseToken(tokenString)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "wrong token",
			})
			c.Abort()
			return

		}
		if claimsmap, ok := claims.(map[string]string); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"userID": claimsmap["userID"],
			})
			c.Next()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "wrong token",
			})
			c.Abort()
			return
		}
	}
}
