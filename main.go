package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go_postgresql/controller"
	"os"
)

func main(){
	router := gin.Default()
	controller.Validate = validator.New()

	// Registering custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("above_age", controller.AboveAge)
		if err != nil {
			return
		}
	}

	router.GET("", controller.Home)

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
