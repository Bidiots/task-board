package controller

import (
	"database/sql"
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"../model"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	db *sql.DB
}

func New(db *sql.DB) *TaskController {
	return &TaskController{
		db: db,
	}
}

func (t *TaskController) RegisterRouter(r gin.IRouter) {
	err := model.CreateTaskTable(t.db)
	if err != nil {
		log.Fatal(err)
	}
	err = model.CreateReceiveTable(t.db)
	if err != nil {
		log.Fatal(err)
	}
	r.POST("/post", t.publish)
	r.POST("/delete", t.deleteByID)

	r.POST("/info/id", t.infoByID)
	r.POST("/info/all", t.infoAll)
	r.GET("/info/download", t.infoAllCsv)
	r.POST("/info/descripty", t.updateDescription)

	r.POST("/info/reveiver", t.infoReveiverBytask)
	r.POST("/user/task", t.infoTaskByuser)
	r.POST("/task/receive", t.receive)
	r.POST("/task/cancle", t.deleteReceive)
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
	_, err = model.InsertTask(t.db, task.Name, task.Description, task.CreateTime, task.Poster)
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

	err = model.DeleteByID(t.db, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
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

	task, err := model.InfoByID(t.db, req.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "task": task})
}

func (t *TaskController) infoAll(c *gin.Context) {
	ban, err := model.InfoAllTask(t.db)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "ban": ban})
}

func (t *TaskController) infoAllCsv(c *gin.Context) {
	tasks, err := model.InfoAllTask(t.db)
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

	poster, err := model.InfoPosterNameByID(t.db, req.ID)
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

	err = model.UpdateDescriptionByID(t.db, req.ID, req.Description)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
