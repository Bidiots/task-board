package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"task-board/task/model"

	"github.com/gin-gonic/gin"
)

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

func (t *TaskController) updateReceiver(c *gin.Context) {
	var req struct {
		userID int
		taskID int
	}
	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest})
		return
	}
	receiver, err := model.InfoReceiverByID(t.db, t.TableName, req.taskID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}
	receiver = fmt.Sprintf(receiver + strconv.Itoa(req.userID))
	err = model.UpdateReceiver(t.db, t.TableName, receiver, req.taskID)

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadGateway, gin.H{"status": http.StatusBadGateway})
		return
	}

}
