package main

import (
	"0chaos.eu/todox/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/tasks", routes.GetTasks)
	r.GET("/tasks/:id", routes.GetTask)
	r.POST("/tasks", routes.PostTask)
	r.DELETE("/tasks/:id", routes.DeleteTask)
	r.Run("localhost:8080")
}
