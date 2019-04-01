package controller

import (
	"database/sql"
	"log"
	"net/http"
	"task-board/task/model"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	db        *sql.DB
	TableName string
}

func New(db *sql.DB, tableName string) *TaskController {
	return &TaskController{
		db:        db,
		TableName: tableName,
	}
}

func (t *TaskController) RegisterRouter(r gin.IRouter) {
	if r == nil {
		log.Fatal("[InitRouter]: server is nil")
	}
	err := model.CreateTable(t.db, t.TableName)
	if err != nil {
		log.Fatal(err)
	}
	r.POST("/post", t.publish)
	r.POST("/delete", t.deleteByID)
	r.POST("/info/id", t.infoByID)
	r.POST("/info/all", t.infoAll)
	r.GET("/info/download", t.infoAllCsv)
	r.POST("/info/descripty", t.updateDescription)
	r.GET("/user/tasks", t.infoUserTask)
	r.POST("/task/accept", t.updateReceiver)
}

func (t *TaskController) publish(c *gin.Context) {
	task := model.Task{}
	err := c.BindJSON(&task)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	if len(task.Name) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "task name can't be empty"})
		return
	}
	if len(task.Description) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "msg": "task descripty can't be empty"})
		return
	}
	task.CreateTime = time.Now()
	_, err = model.InsertTask(t.db, t.TableName, task.Name, task.Description, task.CreateTime, task.Poster)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (t *TaskController) deleteByID(c *gin.Context) {
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
	/*
		handler, err := model.InfoPosterNameByID(t.db, t.TableName, req.ID)

		token, err := c.Cookie("token")
		if err != nil {
			token = "NotSet"
			c.SetCookie("token", "", 3600, "/", "localhost", false, true)
			c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
			return
		}
		if claims, ok := jwt.ParseToken(token); ok {
			if claimsmap, ok := claims.(map[string]string); ok {
				if claimsmap["name"] != handler && claimsmap["type"] != "admin" {
					c.Error(err)
					c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "Without permission"})
					return
				}
			}
	*/
	err = model.DeleteByID(t.db, t.TableName, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	//}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
