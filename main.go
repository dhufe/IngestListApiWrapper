package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	VERSION          = "0.2"
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

	cfgPath, err := ParseFlags()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	router.GET("api", cfg.getDefaultResponse)
	router.POST("/api/upload", cfg.uploadFile)
	router.POST("/api/identify", cfg.identifyFile)
	router.Run(":8080")
}
