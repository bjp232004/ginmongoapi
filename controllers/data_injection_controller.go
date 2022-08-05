package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/blugelabs/bluge"
	"github.com/gin-gonic/gin"
)

var bludgePath = "/home/ubuntu/projects/ginmongoapi/blugedump/"
var folderRootPath = "/home/ubuntu/projects/ginmongoapi/blugedata/"

type userObject struct {
	ID   string
	Name string
}

type Form struct {
	Files []*multipart.FileHeader `form:"files" binding:"required"`
}

func ProcessFiles() {
	go func() {
		fmt.Println("Inside ProcessFiles")
		config := bluge.DefaultConfig(bludgePath)
		writer, err := bluge.OpenWriter(config)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		defer writer.Close()

		files, err := ioutil.ReadDir(folderRootPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			fileName := folderRootPath + f.Name()
			fmt.Println(fileName)

			doc := bluge.NewDocument(fileName)
			docs, err := ioutil.ReadFile(fileName)

			if err != nil {
				fmt.Printf("error while readfile: %v", err)
				os.Exit(1)
			}

			var users []userObject

			err = json.Unmarshal(docs, &users)
			if err != nil {
				fmt.Printf("error while unmarshal: %v", err)
				os.Exit(1)
			}

			for k := range users {
				doc.AddField(bluge.NewTextField("id", users[k].ID))
				doc.AddField(bluge.NewTextField("name", users[k].Name))
				err = writer.Update(doc.ID(), doc)
			}

			if err != nil {
				fmt.Printf("error updating document: %v", err)
				os.Exit(1)
			}
		}
	}()
	return
}

func CreateZincIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Print("Create Index Started")
		// config := bluge.DefaultConfig(bludgePath)
		// writer, err := bluge.OpenWriter(config)
		// if err != nil {
		// 	fmt.Printf("%v\n", err)
		// }
		// defer writer.Close()

		// files, err := ioutil.ReadDir(folderRootPath)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// for _, f := range files {
		// 	fileName := folderRootPath + f.Name()
		// 	fmt.Println(fileName)

		// 	doc := bluge.NewDocument(fileName)
		// 	docs, err := ioutil.ReadFile(fileName)

		// 	if err != nil {
		// 		fmt.Printf("error while readfile: %v", err)
		// 		os.Exit(1)
		// 	}

		// 	var users []userObject

		// 	err = json.Unmarshal(docs, &users)
		// 	if err != nil {
		// 		fmt.Printf("error while unmarshal: %v", err)
		// 		os.Exit(1)
		// 	}

		// 	for k := range users {
		// 		doc.AddField(bluge.NewTextField("id", users[k].ID))
		// 		doc.AddField(bluge.NewTextField("name", users[k].Name))
		// 		err = writer.Update(doc.ID(), doc)
		// 	}

		// 	if err != nil {
		// 		fmt.Printf("error updating document: %v", err)
		// 		os.Exit(1)
		// 	}
		// }
		ProcessFiles()

		c.JSON(http.StatusOK, gin.H{"status": "Index has been created."})
	}
}

func UploadFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Print("File upload started")

		form, _ := c.MultipartForm()

		// files := form.File["files"]

		// for key, value := range c.Request.PostForm {
		// 	log.Printf("Filename from: %v = %v \n", key, value)
		// }

		for _, file := range form.File["zincsearchdata.json"] {
			// log.Println(file.Filename)

			// Upload the file to specific dst.
			c.SaveUploadedFile(file, folderRootPath+file.Filename)
		}

		ProcessFiles()

		// Using `ShouldBind`
		// --------------------
		// var form Form
		// if err := c.ShouldBind(&form); err != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 	return
		// }

		// for _, formFile := range form.Files {

		// 	// Get raw file bytes - no reader method
		// 	// openedFile, _ := formFile.Open()
		// 	// file, _ := ioutil.ReadAll(openedFile)

		// 	// Upload to disk
		// 	// `formFile` has io.reader method
		// 	c.SaveUploadedFile(formFile, folderRootPath)

		// }
		// --------------------

		// Using `MultipartForm`
		// --------------------
		// form, _ := c.MultipartForm()
		// formFiles, _ := form.File["file"]
		// fmt.Printf("%v", form.File["file"])

		// for _, formFile := range formFiles {

		// 	// Get raw file bytes - no reader method
		// 	// openedFile, _ := formFile.Open()
		// 	// file, _ := ioutil.ReadAll(openedFile)
		// 	// fmt.Printf("%v", file)
		// 	// Upload to disk
		// 	// `formFile` has io.reader method
		// 	c.SaveUploadedFile(formFile, folderRootPath)

		// }
		// --------------------

		// fmt.Printf("%s", c.Request)

		// reader, _ := c.Request.MultipartReader()
		// for {
		// 	part, err := reader.NextPart()
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))

		// 	if err != nil {
		// 		break
		// 	}
		// 	// filename := params["filename"]
		// 	fmt.Printf("%s", params)
		// }

		// file, header, err := c.Request.FormFile("files")
		// // fmt.Printf("%s", file)
		// if err != nil {
		// 	c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		// 	return
		// }
		// filename := header.Filename
		// out, err := os.Create(folderRootPath + filename)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer out.Close()
		// _, err = io.Copy(out, file)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
	}
}

func SearchData() gin.HandlerFunc {
	return func(c *gin.Context) {
		matchFlag := "Not Match"
		searchColumn := c.Param("field")
		searchTerm := c.Param("value")
		fmt.Printf("Search Term: %s", searchTerm)
		config := bluge.DefaultConfig(bludgePath)
		writer, err := bluge.OpenWriter(config)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		defer writer.Close()

		//QUery Index
		reader, err := writer.Reader()
		if err != nil {
			fmt.Printf("error getting index reader: %v", err)
		}
		defer reader.Close()

		query := bluge.NewMatchQuery(searchTerm).SetField(searchColumn)

		request := bluge.NewTopNSearch(100, query).
			WithStandardAggregations()
		documentMatchIterator, err := reader.Search(context.Background(), request)
		if err != nil {
			fmt.Printf("error executing search: %v", err)
		}
		match, err := documentMatchIterator.Next()
		for err == nil && match != nil {
			err = match.VisitStoredFields(func(field string, value []byte) bool {
				if field == "_id" {
					matchFlag = "Matched"
					fmt.Printf("match found: %s\n", string(value))
				}
				return true
			})
			if err != nil {
				fmt.Printf("error loading stored fields: %v", err)
			}
			match, err = documentMatchIterator.Next()
		}
		if err != nil {
			fmt.Printf("error iterator document matches: %v", err)
		}
		c.JSON(http.StatusOK, gin.H{"match found end": matchFlag})
	}
}
