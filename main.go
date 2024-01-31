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
}

func main() {
	defer initializers.CloseDBConnection(initializers.DB)

	r := gin.Default()

	r.POST("/create", controllers.PeopleCreate)
	r.PUT("/update/:id", controllers.PeopleUpdate)
	r.DELETE("/delete/:id", controllers.PeopleDelete)
	r.GET("/get", controllers.PeopleGet)

	r.Run()
}
