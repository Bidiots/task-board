package controller

import (
	"net/http"
	"task-board/permission/model"

	"github.com/gin-gonic/gin"
)

func (c *Controller) createRole(ctx *gin.Context) {
	var (
		role struct {
			Name  string `json:"name"`
			Intro string `json:"intro"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.CreateRole(c.db, role.Name, role.Intro)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) modifyRole(ctx *gin.Context) {
	var (
		role struct {
			RoleID int    `json:"role_id"`
			Name   string `json:"name"`
			Intro  string `json:"intro"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.ModifyRole(c.db, role.RoleID, role.Name, role.Intro)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) roleList(ctx *gin.Context) {
	result, err := model.InfoRoleList(c.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "RoleList": result})
}

func (c *Controller) getRoleByID(ctx *gin.Context) {
	var (
		role struct {
			RoleID int `json:"role_id"`
		}
	)

	err := ctx.ShouldBind(&role)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	result, err := model.GetRoleByID(c.db, role.RoleID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "RoleByID": result})
}
