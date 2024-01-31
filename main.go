package main

import (
	"example/m/controllers"
	"example/m/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.InitLogrus()

	defer initializers.CloseDBConnection(initializers.DB)
}

func main() {
	r := gin.Default()

	r.POST("/create", controllers.PeopleCreate)
	r.PUT("/update", controllers.PeopleUpdate)
	r.DELETE("/delete", controllers.PeopleDelete)
	r.GET("/get", controllers.PeopleGet)

	r.Run()
}
