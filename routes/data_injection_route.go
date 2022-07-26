package routes

import (
	"ginmongoapi/controllers" //add this

	"github.com/gin-gonic/gin"
)

func DataInjectionRoute(router *gin.Engine) {
	router.POST("/upload", controllers.UploadFile())
	router.GET("/file_upload", controllers.FileUpload())
}
