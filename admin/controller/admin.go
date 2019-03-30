package controller

import (
	"TEST/admin/model"
	"TEST/jwt"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	db        *sql.DB
	tableName string
}

func NewAdminController(db *sql.DB, tableName string) *AdminController {
	return &AdminController{
		db:        db,
		tableName: tableName,
	}
}

func (a *AdminController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := model.CreateTable(a.db, a.tableName)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/register", a.register)
	r.POST("/delete/:id", a.deleteByID)
	r.POST("/info/:id", a.infoByID)
	r.POST("/login", a.login)
}
func (a *AdminController) infoByID(c *gin.Context) {
	var (
		req struct {
			ID int `json:"id"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	ban, err := model.InfoByID(a.db, a.tableName, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ban": ban})
}
func (a *AdminController) deleteByID(c *gin.Context) {
	var (
		req struct {
			ID int `json:"id"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.DeleteByID(a.db, a.tableName, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
func (a *AdminController) register(c *gin.Context) {
	var user model.Admin
	err := c.BindJSON(&user)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	if user.Name != "" && user.Password != "" {
		_, err = model.InsertAdmin(a.db, a.tableName, user.Name, user.Password)
		if err != nil {
			log.Fatal(err)
		}
		m := make(map[string]interface{}, 3)
		m["name"] = user.Name
		m["string"] = user.Password
		m["type"] = "admin"
		tokenString := jwt.CreateToken(m)
		c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

		c.String(http.StatusOK, "注册成功")
	} else {
		log.Println("用户名不能为空")
	}

}
func (a *AdminController) login(c *gin.Context) {
	user := &model.Admin{}
	err := c.Bind(user)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	if user.Name != "" && user.Password != "" {
		password, err := model.InfoPasswordByName(a.db, a.tableName, user.Name)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
			return
		}
		if password == user.Password {
			m := make(map[string]interface{}, 3)
			m["name"] = user.Name
			m["string"] = user.Password
			m["type"] = "admin"
			tokenString := jwt.CreateToken(m)
			c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

		} else {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "密码错误"})
		}

		c.String(http.StatusOK, "登陆成功")
	}
}
