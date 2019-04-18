package controller

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"../../config"
	"../../user/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//UserController -
type UserController struct {
	db        *sql.DB
	tableName string
}

//New -
func New(db *sql.DB, tableName string) *UserController {
	return &UserController{
		db:        db,
		tableName: tableName,
	}
}

//RegisterRouter -
func (u *UserController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}

	err := model.CreateTable(u.db, u.tableName)
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/delete/:id", u.deleteByID)
	r.GET("/info/:id", u.infoByID)
}

//Register -
func (u *UserController) Register(c *gin.Context) {
	var user1 model.User
	err := c.BindJSON(&user1)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	if user1.Name != "" && user1.Password != "" {
		userID, err := model.InsertUser(u.db, u.tableName, user1.Name, user1.Password)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
			return
		}

		m := make(map[string]interface{}, 2)
		m["userID"] = userID
		m["name"] = user1.Name

		token := jwt.New(jwt.SigningMethodHS256)
		claims := make(jwt.MapClaims)

		for index, val := range m {
			claims[index] = val
		}

		token.Claims = claims

		tokenString, _ := token.SignedString([]byte(config.Key))

		c.Header("Authorization", tokenString)

		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "register ok"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})

	}

}

//infoByID -
func (u *UserController) infoByID(c *gin.Context) {
	ID := c.Param("id")

	IDint, err := strconv.Atoi(ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	user, err := model.InfoByID(u.db, u.tableName, IDint)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "user": user})
}

//deleteByID -
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

//Login -
func (u *UserController) Login(c *gin.Context) {
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
			m := make(map[string]interface{}, 2)
			m["name"] = user.Name
			m["userID"] = user.ID

			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)

			for index, val := range m {
				claims[index] = val
			}

			token.Claims = claims
			tokenString, _ := token.SignedString([]byte(config.Key))

			c.Header("Authorization", tokenString)

			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "login success"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "password is wrong"})
		}

	}
}
