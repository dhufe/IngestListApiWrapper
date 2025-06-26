package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fileUploadResponse struct {
	DurationInMs int64  `json:"durationInMs"`
	FilePath     string `json:"filePath"`
}

type fileIdentifyRequest struct {
	FilePath string `json:"filePath" binding:"required"`
}

func getDefaultResponse(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf(DEFAULT_RESPONSE, VERSION))
}

func identifyFile(c *gin.Context) {
	var filePath fileIdentifyRequest

	if err := c.ShouldBind(&filePath); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "path is missing",
		})
		return
	}

	fmt.Printf("Processing : %s", filePath.FilePath)

	if _, err := os.Stat(filePath.FilePath); err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "file does not exist",
		})
		return
	}

	con, err := net.Dial("tcp", DEFAULT_IL_SERVER)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "unable to open connection to socket",
		})
		return
	}
	defer con.Close()

	// Send data to socket
	_, err = con.Write([]byte(filePath.FilePath))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "unable to write to socket",
		})
		return
	}
}

func uploadFile(c *gin.Context) {
	s := time.Now()

	file, err := c.FormFile("file")
	// no file received
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "no file received",
		})
		return
	}
	// generate unique file name for storing
	filename := uuid.New().String() + "_" + file.Filename
	fileStorePath := filepath.Join(FILE_STORE_PATH, filename)
	err = c.SaveUploadedFile(file, fileStorePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "unable to save file",
		})
		return
	}

	response := fileUploadResponse{
		DurationInMs: time.Since(s).Milliseconds(),
		FilePath:     fileStorePath,
	}

	c.JSON(http.StatusOK, response)
	//	defer os.Remove(fileStorePath)
}
