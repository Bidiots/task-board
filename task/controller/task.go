package controller

import (
	"TEST/jwt"
	"TEST/task/model"
	"database/sql"
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
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
	r.GET("/id/tasks", t.showUserTask)
}
func (t *TaskController) infoByID(c *gin.Context) {
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
	task, err := model.InfoByID(t.db, t.TableName, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "task": task})
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
		err = model.DeleteByID(t.db, t.TableName, req.ID)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
			return
		}

	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
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
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
func (t *TaskController) infoAll(c *gin.Context) {
	ban, err := model.InfoAllTask(t.db, t.TableName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ban": ban})
}
func (t *TaskController) infoAllCsv(c *gin.Context) {
	tasks, err := model.InfoAllTask(t.db, t.TableName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	file, err := os.Create("task.csv")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	defer file.Close()

	csv.NewReader(file)
	for _, task := range tasks {
		_, err := file.WriteString(strconv.Itoa(task.ID) + ",")
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.WriteString(task.Name + ",")
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.WriteString(task.Description + ",")
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.WriteString(task.Receiver + ",")
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.WriteString(task.Poster + ",")
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.WriteString("\n")
		if err != nil {
			log.Fatal(err)
		}

	}
}
func (t *TaskController) updateDescription(c *gin.Context) {
	var (
		req struct {
			Description string `json:"description"`
			ID          int    `json:"id"`
			User        string `json:"user"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	poster, err := model.InfoPosterNameByID(t.db, t.TableName, req.ID)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	if ok, _ := regexp.MatchString(req.User, poster); !ok {
		c.Error(err)
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "Without Permission"})
		return
	}
	err = model.UpdateDescriptionByID(t.db, t.TableName, req.ID, req.Description)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})

}
func (t *TaskController) showUserTask(c *gin.Context) {
	var req struct {
		userName string
	}
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	tasks, err := model.InfoByReceiver(t.db, t.TableName, req.userName)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": http.StatusOK, "tasks": tasks})
}
