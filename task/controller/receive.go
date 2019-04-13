package controller

import (
	"net/http"

	"../model"
	"github.com/gin-gonic/gin"
)

func (t *TaskController) receive(c *gin.Context) {
	var req struct {
		userID int `json:"userID"`
		taskID int `json:"taskID"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}

	err = model.InsertReceiverInfo(t.db, req.userID, req.taskID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}

func (t *TaskController) infoReveiverBytask(c *gin.Context) {
	var req struct {
		taskID int `json:"taskID"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	tasksID, err := model.QueryUserIDByTaskID(t.db, req.taskID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	var tasks []*model.Task
	for _, taskID := range *tasksID {
		task, err := model.InfoByID(t.db, taskID)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
			return
		}
		tasks = append(tasks, task)
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "tasks": tasks})
}

func (t *TaskController) infoTaskByuser(c *gin.Context) {
	var req struct {
		userID int `json:"userID"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	usersID, err := model.QueryTaskByUserID(t.db, req.userID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "usersID": usersID})
}

func (t *TaskController) deleteReceive(c *gin.Context) {
	var req struct {
		userID int `json:"userID"`
		taskID int `json:"taskID"`
	}
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	err = model.DeleteReceiveInfo(t.db, req.userID, req.taskID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK})
}
