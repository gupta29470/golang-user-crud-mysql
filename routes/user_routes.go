package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gupta29470/golang_sql_crud_without_orm/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users", controllers.CreateUser())
	incomingRoutes.GET("/users", controllers.GetAllUsers())
	incomingRoutes.GET("/users/:id", controllers.GetUser())
	incomingRoutes.PUT("/users/:id", controllers.UpdateUser())
	incomingRoutes.DELETE("users/:id", controllers.DeleteUser())

}
