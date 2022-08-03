package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go_postgresql/controller"
	"os"
	"time"
)

func main(){
	fmt.Println(time.Now().Add(-time.Hour * 140160).Unix())

	router := gin.Default()
	controller.Validate = validator.New()

	apiRoutes := router.Group("/api")

	apiRoutes.GET("/all", controller.GetAllUsers)
	apiRoutes.POST("/register", controller.RegisterUser)
	apiRoutes.POST("/login", controller.LoginUser)
	apiRoutes.PATCH("/update/:id", controller.UpdateDetails)
	apiRoutes.DELETE("/delete/:id", controller.DeleteUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		return
	}
}
