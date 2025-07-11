package storage

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"dhufe/ingestlistapiwrapper/models"
)

const (
	VERSION         = "0.3"
	DefaultResponse = "IngestList-Wrapper version %s is running"
)

type Repository struct {
	DataBase  *gorm.DB
	Config    *Config
	Scheduler *gocron.Scheduler
}

func (r *Repository) getDefaultResponse(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf(DefaultResponse, VERSION))
}

func (r *Repository) CreateRoutes(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/", r.getDefaultResponse)
	router.GET("/", r.getDefaultResponse)
	api.POST("/upload", r.CreateJob)
	api.GET("/jobs", r.GetJobs)

	api.GET("/job/:id", r.GetJobByID)
	api.DELETE("/job/:id", r.DeleteJob)
}

func (r *Repository) CreateJob(c *gin.Context) {
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
	fileStorePath := filepath.Join(r.Config.FileStorePath, filename)
	err = c.SaveUploadedFile(file, fileStorePath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "unable to save file",
		})
		return
	}

	status := "New"

	job := models.Jobs{
		FilePath: &fileStorePath,
		Status:   &status,
		Created:  s,
		Result:   nil,
	}

	res := r.DataBase.Create(&job)
	if res.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "could not create job",
		})
	}

	response := fileUploadResponse{
		Data: job,
	}

	c.JSON(http.StatusOK, response)
}

func (r *Repository) DummyFunc() {
	jobModels := &[]models.Jobs{}

	err := r.DataBase.Where("status not LIKE ?", "Finished").Find(&jobModels).Error
	if err == nil {
		fmt.Printf("%s Scheduler executed.\n", time.Now())
		fmt.Printf("Currently there are %d waiting tasks.\n", len(*jobModels))
	}
}

func (r *Repository) CleaningDatabase() {
	jobModels := &[]models.Jobs{}
	result := r.DataBase.Where("created < CURRENT_DATE - INTERVAL '7 days' AND status LIKE ?", "Finished").
		Find(&jobModels)
	if result.Error == nil {
		fmt.Printf("%s database cleanup executed.\n", time.Now())
		fmt.Printf("This job would wipe %d jobs.\n", len(*jobModels))
	}
}

func (r *Repository) ProcessEntry() {
	jobModel := &models.Jobs{}

	err := r.DataBase.Find(&jobModel, "status NOT LIKE ?", "Finished").Order("id DESC")
	if err.RowsAffected != 0 {

		//	if err == nil {
		fmt.Printf("Processing %s.\n", *jobModel.FilePath)

		if _, err := os.Stat(*jobModel.FilePath); err != nil {
			fmt.Println(err.Error())
			return
		}

		con, err := net.Dial("tcp", r.Config.IngestListServer)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer con.Close()

		// Send data to socket
		n, err := con.Write([]byte(*jobModel.FilePath + "\n"))
		fmt.Printf("Written %d bytes to socket.\n", n)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		msg := make([]byte, 8192)
		n, err = io.ReadFull(con, msg)
		fmt.Printf("Read %d bytes from socket.\n", n)

		if err != nil && err != io.ErrUnexpectedEOF {
			fmt.Println(err.Error())
			return
		}

		var result string
		result = string(msg[:n])
		status := "Finished"

		j := models.Job{
			FilePath: jobModel.FilePath,
			Status:   &status,
			Created:  jobModel.Created,
			Result:   &result,
		}

		err = r.DataBase.Model(&jobModel).Updates(j).Error
		if err != nil {
			fmt.Printf("Can not update entry")
			return
		}
		// defer os.Remove(*jobModel.FilePath)
	}
}

func (r *Repository) AddJob() {
	task := gocron.NewTask(
		r.ProcessEntry,
	)

	job := gocron.CronJob(
		r.Config.SchedulerTab,
		false,
	)
	j, err := (*r.Scheduler).NewJob(job, task)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Tick-Job id = %d.\n", j.ID())

	cleanUpTask := gocron.NewTask(
		r.CleaningDatabase,
	)

	cleanUpJob := gocron.CronJob(
		"0 0 * * *",
		false,
	)

	j, err = (*r.Scheduler).NewJob(cleanUpJob, cleanUpTask)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("CleanUp-Job id = %d.\n", j.ID())
}

func (r *Repository) DeleteJob(c *gin.Context) {
	jobModel := &models.Jobs{}

	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id cannot be empty",
		})
		return
	}

	err := r.DataBase.Where("id = ?", id).First(jobModel).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "could not get the job",
		})
		return
	}

	err = os.Remove(*jobModel.FilePath)
	if err != nil {
		return
	}
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	//		"message": "could not delete file in filesystem.",
	//	})
	//	return
	//}

	err = r.DataBase.Delete(models.Jobs{}, id).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "could not delete job",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "job delete successfully", // cast it to string before showing
	})
}

func (r *Repository) GetJobs(c *gin.Context) {
	jobModels := &[]models.Jobs{}

	err := r.DataBase.Find(jobModels).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "could not get jobs",
		})
		return
	}

	response := jobsResponse{
		Message: "success",
		Data:    jobModels,
	}

	c.JSON(http.StatusOK, response)
}

func (r *Repository) GetJobByID(c *gin.Context) {
	id := c.Param("id")
	jobModel := &models.Jobs{}

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id cannot be empty",
		})
		return
	}
	err := r.DataBase.Where("id = ?", id).First(jobModel).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "could not get the job",
		})
		return
	}

	response := jobByIdResponse{
		Message: "success",
		Data:    *jobModel,
	}

	c.JSON(http.StatusOK, response)
}
