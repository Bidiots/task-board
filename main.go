package main

import (
	"database/sql"
	"net/http"
	permission "task-board/permission/controller"
	task "task-board/task/controller"
	user "task-board/user/controller"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
	dbConn, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/test1?charset=utf8&parseTime=true")
	if err != nil {
		panic(err)
	}
	UserCon := user.New(dbConn, "user1")
	TaskCon := task.New(dbConn, "task")
	Permission := permission.New(dbConn)
	router.POST("/user/login", UserCon.Login)
	router.Use(UserCon.MiddleWareJWT())
	router.Use(Permission.CheckPermission())
	Permission.RegisterRouter(router.Group("/permission"))
	UserCon.RegisterRouter(router.Group("/user"))
	TaskCon.RegisterRouter(router.Group("/task"))
	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
