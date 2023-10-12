package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gupta29470/golang_sql_crud_without_orm/database"
	"github.com/gupta29470/golang_sql_crud_without_orm/routes"
)

func init() {
	database.InitDB()
}

func main() {
	router := gin.New()

	routes.UserRoutes(router)
	router.Run(":9000")
}
