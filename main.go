package main

import (
	admin "TEST/admin/controller"
	task "TEST/task/controller"
	user "TEST/user/controller"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()
	dbConn, err := sql.Open("mysql", "root:123@tcp(127.0.0.1:3306)/test1?charset=utf8")
	if err != nil {
		panic(err)
	}
	UserCon := user.New(dbConn, "user1")
	TaskCon := task.New(dbConn, "task")
	AdminCon := admin.New(dbConn, "admin")
	AdminCon.RegisterRouter(router.Group("/admin"))
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
