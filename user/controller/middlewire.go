package controller

import (
	"net/http"

	"../../config"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

//CheckJWT -
func (u *UserController) CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {
				return []byte(config.Key), nil
			})

		if err == nil {
			if token.Valid {
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized})
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized})
		c.Abort()
		return

	}
}

/*
func (u *UserController) MiddleWareJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
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
*/
