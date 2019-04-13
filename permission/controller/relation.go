package controller

import (
	"net/http"

	"../../permission/model"

	"github.com/gin-gonic/gin"
)

func (c *Controller) addRelation(ctx *gin.Context) {
	var (
		relation struct {
			AdminID int `json:"admin_id"`
			RoleID  int `json:"role_id"`
		}
	)

	err := ctx.ShouldBind(&relation)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.InsertRelation(c.db, relation.AdminID, relation.RoleID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) removeRelation(ctx *gin.Context) {
	var (
		relation struct {
			AdminID int `json:"admin_id"`
			RoleID  int `json:"role_id"`
		}
	)

	err := ctx.ShouldBind(&relation)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.DeleteRelation(c.db, relation.AdminID, relation.RoleID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
