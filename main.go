package main

import (
	"ginmongoapi/configs"
	"ginmongoapi/routes" //add this

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(router)          //add this
	routes.DataInjectionRoute(router) //add this

	router.Run("localhost:8080")
}
