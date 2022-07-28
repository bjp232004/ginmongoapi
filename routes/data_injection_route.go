package routes

import (
	"ginmongoapi/controllers" //add this

	"github.com/gin-gonic/gin"
)

func DataInjectionRoute(router *gin.Engine) {
	router.POST("/upload", controllers.UploadFile())
	router.GET("/create_index", controllers.CreateZincIndex())
	router.GET("/search/:field/:value", controllers.SearchData())
}
