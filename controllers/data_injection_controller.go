package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func UploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
			return
		}
		filename := header.Filename
		out, err := os.Create("public/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		filepath := "http://localhost:8080/file/" + filename
		c.JSON(http.StatusOK, gin.H{"filepath": filepath})
	}
}

func FileUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "select_file.html", gin.H{})
	}
}
