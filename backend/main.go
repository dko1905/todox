package main

import (
	"log"
	"os"
	"time"

	"0chaos.eu/todox/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: false,
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		MaxAge:           360 * time.Second,
	}))

	r.GET("/tasks", routes.GetTasks)
	r.GET("/tasks/:id", routes.GetTask)
	r.POST("/tasks", routes.PostTask)
	r.DELETE("/tasks/:id", routes.DeleteTask)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on", port)
	log.Fatal(r.Run(":" + port))
}
