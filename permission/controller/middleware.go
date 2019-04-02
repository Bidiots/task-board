package controller

import (
	"net/http"
	"strconv"
	"task-board/jwt"
	"task-board/permission/model"

	"github.com/gin-gonic/gin"
)

func (c *Controller) CheckPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		check := false
		tokenString := ctx.Request.Header.Get("token")
		if tokenString == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "No Token",
			})
			ctx.Abort()
			return
		}
		claims, ok := jwt.ParseToken(tokenString)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "Wrong Token",
			})
			ctx.Abort()
			return

		}
		var IDString string
		claimsmap, ok := claims.(map[string]string)
		if ok {
			IDString = claimsmap["userID"]
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"msg":    "Wrong Token",
			})
			ctx.Abort()
			return
		}
		userID, err := strconv.Atoi(IDString)
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
