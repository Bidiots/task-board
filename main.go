package main

import (
	"database/sql"
	"net/http"
	"time"

	permission "./permission/controller"
	task "./task/controller"
	user "./user/controller"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()

	dbConn, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:8806)/test1?charset=utf8&parseTime=true")
	if err != nil {
		panic(err)
	}

	UserCon := user.New(dbConn, "user1")
	TaskCon := task.New(dbConn)
	Permission := permission.New(dbConn)

	router.POST("/user/register", UserCon.Register)
	router.POST("/user/login", UserCon.Login)

	router.Use(UserCon.CheckJWT())
	//router.Use(Permission.CheckPermission())

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
