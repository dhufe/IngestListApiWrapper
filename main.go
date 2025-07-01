package main

import (
	"dhufe/ingestlistapiwrapper/models"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"dhufe/ingestlistapiwrapper/storage"
)

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 3000 << 20 // 3 GiB
	// Allow cors to integrate Borg in other applications.
	router.ForwardedByClientIP = true
	router.SetTrustedProxies([]string{"*"})
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	// It's important that the cors configuration is used before declaring the routes.
	router.Use(cors.New(corsConfig))

	cfgPath, err := storage.ParseFlags()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	cfg, err := storage.NewConfig(cfgPath)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	db, err := storage.NewDatabase(&cfg.DbConfig)
	if err != nil {
		fmt.Println("could not load the database")
	}

	if db != nil {
		fmt.Printf("database loaded successfully")
		models.MigrateJobs(db)
	}

	r := storage.Repository{
		DataBase:  db,
		Config:    cfg,
		Scheduler: NewScheduler(),
	}

	r.AddJob()

	(*r.Scheduler).Start()

	r.CreateRoutes(router)
	router.Run(":8080")
}
