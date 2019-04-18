package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"../../permission/model"

	"github.com/gin-gonic/gin"
)

//CheckPermission -
func (c *Controller) CheckPermission() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		check := false
		tokenString := ctx.GetHeader("Authorizatio")
		if tokenString == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "Without Token",
			})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("12345"), nil
		})

		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "Wrong Token",
			})
			ctx.Abort()
			return
		}

		IDstring := fmt.Sprint(claims["userID"])
		userID, err := strconv.Atoi(IDstring)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
			ctx.Abort()
		}

		URLL := ctx.Request.URL.Path

		adRole, err := model.GetRoleMap(c.db, userID)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusConflict, gin.H{
				"status": http.StatusConflict,
			})
			ctx.Abort()
			return
		}

		urlRole, err := model.URLPermissionsMap(c.db, URLL)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusFailedDependency, gin.H{
				"status": http.StatusFailedDependency,
			})
			ctx.Abort()
			return
		}

		for urlkey := range urlRole {
			for adkey := range adRole {
				if urlkey == adkey {
					check = true
				}
			}
		}

		if !check {
			ctx.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"msg":    "Without Permission",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
