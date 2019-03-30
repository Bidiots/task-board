package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"task-board/jwt"
	"task-board/user/model"

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
	r.GET("/delete/:id", u.deleteByID)
	r.GET("/info/:id", u.infoByID)
	r.POST("/login", u.login)
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
	ID := c.Param("id")
	IDint, err := strconv.Atoi(ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusOK, gin.H{"status": http.StatusBadRequest})
		return
	}
	err = model.DeleteByID(u.db, u.tableName, IDint)
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
			c.Error(err)
			c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
			return
		}
		m := make(map[string]interface{}, 3)
		m["name"] = user1.Name
		m["string"] = user1.Password
		m["type"] = "user"
		tokenString := jwt.CreateToken(m)
		c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "register ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})

	}

}
func (u *UserController) login(c *gin.Context) {
	user := &model.User{}
	err := c.Bind(user)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	if user.Name != "" && user.Password != "" {
		password, err := model.InfoPasswordByName(u.db, u.tableName, user.Name)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
			return
		}
		if password == user.Password {
			m := make(map[string]interface{}, 3)
			m["name"] = user.Name
			m["string"] = user.Password
			m["type"] = "user"
			tokenString := jwt.CreateToken(m)
			c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)

		} else {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "password is wrong"})
		}

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "login ok"})
	}
}
