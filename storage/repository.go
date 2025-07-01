package storage

import (
	"dhufe/ingestlistapiwrapper/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"time"
)

const (
	VERSION         = "0.3"
	DefaultResponse = "IngestList-Wrapper version %s is running"
)

type Repository struct {
	DataBase *gorm.DB
	Config   *Config
}

func (r *Repository) getDefaultResponse(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf(DefaultResponse, VERSION))
}

func (r *Repository) CreateRoutes(router *gin.Engine) {
	router.GET("api", r.getDefaultResponse)
	router.POST("/api/upload", r.CreateJob)
	router.GET("/api/jobs", r.GetJobs)
	router.GET("/api/job/:id", r.GetJobByID)
	router.DELETE("/api/job/:id", r.DeleteJob)
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

	response := fileUploadResponse{
		DurationInMs: time.Since(s).Milliseconds(),
		FilePath:     fileStorePath,
	}

	var status = "New"
	job := models.Job{
		&fileStorePath,
		&status,
		s,
	}

	err = r.DataBase.Create(&job).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "could not create job",
		})

	}

	c.JSON(http.StatusOK, response)
}

func (r *Repository) DeleteJob(c *gin.Context) {
	jobModel := models.Jobs{}

	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "id cannot be empty",
		})
		return
	}

	err := r.DataBase.Delete(jobModel, id)

	if err.Error != nil {
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

	fmt.Println("the ID is", id)
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
