package controller

import (
	"database/sql"
	"log"
	"net/http"
	"task-board/permission/model"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	db *sql.DB
}

func New(db *sql.DB) *Controller {
	return &Controller{
		db: db,
	}
}

func (c *Controller) RegisterRouter(r gin.IRouter) {
	err := model.CreateRoleTable(c.db)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/addrole", c.createRole)
	r.POST("/modifyrole", c.modifyRole)
	r.POST("/info/all", c.roleList)
	r.POST("/info/id", c.getRoleByID)

	r.POST("/addurl", c.addURLPermission)
	r.POST("/deleteurl", c.removeURLPermission)
	r.POST("/urlgetrole", c.urlPermissions)
	r.POST("/geturl", c.permissions)

	r.POST("/addrelation", c.addRelation)
	r.POST("/removerelation", c.removeRelation)

}

func (c *Controller) addURLPermission(ctx *gin.Context) {
	var (
		url struct {
			URL    string `json:"url"`
			RoleID int    `json:"role_id"`
		}
	)

	err := ctx.ShouldBind(&url)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.InsertURLPermission(c.db, url.RoleID, url.URL)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) removeURLPermission(ctx *gin.Context) {
	var (
		url struct {
			URL    string `json:"url"`
			RoleID int    `json:"role_id"`
		}
	)

	err := ctx.ShouldBind(&url)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.DeleteURLPermission(c.db, url.RoleID, url.URL)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (c *Controller) urlPermissions(ctx *gin.Context) {
	var (
		url struct {
			URL string `json:"url"`
		}
	)

	err := ctx.ShouldBind(&url)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	result, err := model.InfoURLPermissions(c.db, url.URL)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "URLPermissions": result})
}

func (c *Controller) permissions(ctx *gin.Context) {
	result, err := model.InfoPermissions(c.db)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "Permissions": result})
}
