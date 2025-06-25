package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	FILE_STORE_PATH  = "/data"
	VERSION          = "0.1"
	DEFAULT_RESPONSE = "IngestList-Wrapper version %s is running"
)

type fileIndentify struct {
	DurationInMs int64 `json:"durationInMs"`
}

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 3000 << 20 // 3 GiB
	// Allow cors to integrate Borg in other applications.
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"*"})
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	corsConfig.AllowMethods = []string{"GET", "POST"}
	// It's important that the cors configuration is used before declaring the routes.
	router.Use(cors.New(corsConfig))
	router.GET("api", getDefaultResponse)

	router.POST("/api/upload", uploadFile)
	router.GET("/api/identify", identifyFile)

	router.Run("localhost:8080")
}
