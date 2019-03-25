package controller

import (
	"TEST/jwt"
	"TEST/user/model"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	db        *sql.DB
	tableName string
}

func New(db *sql.DB, tableName string) *UserController {
	return &UserController{
		db:        db,
		tableName: tableName,
	}
}
func (u *UserController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := model.CreateTable(u.db, u.tableName)
	if err != nil {
		log.Fatal(err)
	}

	r.POST("/register", u.register)
	r.POST("/delete", u.deleteByID)
	r.POST("/info/id", u.infoByID)
}
func (u *UserController) infoByID(c *gin.Context) {
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

	ban, err := model.InfoByID(u.db, u.tableName, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ban": ban})
}
func (u *UserController) deleteByID(c *gin.Context) {
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

	err = model.DeleteByID(u.db, u.tableName, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
func (u *UserController) register(c *gin.Context) {
	var user1 model.User
	err := c.BindJSON(&user1)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	if user1.Name != "" && user1.Password != "" {
		_, err = model.InsertUser(u.db, u.tableName, user1.Name, user1.Password)
		if err != nil {
			log.Fatal(err)
		}
		m := make(map[string]interface{}, 2)
		m["name"] = user1.Name
		m["string"] = user1.Password
		tokenString := jwt.CreateToken(m)
		c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

		c.String(http.StatusOK, "注册成功")
	} else {
		log.Println("用户名不能为空")
	}

}
