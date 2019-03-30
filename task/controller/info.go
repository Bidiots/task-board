package controller

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"
	"task-board/task/model"

	"github.com/gin-gonic/gin"
)

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

func (t *TaskController) infoUserTask(c *gin.Context) {
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
