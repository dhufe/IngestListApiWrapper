package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hufschlager.net/IngestListApiWrapper/internal/application/services"
	"hufschlager.net/IngestListApiWrapper/internal/domain/tasks/models"
)

type TaskHandler struct {
	service *services.TaskService
}

const (
	VERSION         = "0.4"
	DefaultResponse = "IngestList-Wrapper version %s is running"
	FileUploadType  = "multipart/form-data"
	FileApiType     = "application/json"
)

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var request struct {
		FileName string `json:"filename" binding:"required"`
		Type     string `json:"type" binding:"required"`
	}

	switch c.ContentType() {
	case FileUploadType:
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var tmp struct {
			Type     string                `form:"type" binding:"required"`
			fileName *multipart.FileHeader `form:"file" binding:"required"`
		}

		if err = c.Bind(&tmp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fileName := uuid.New().String() + "_" + file.Filename
		request.FileName = filepath.Join(h.service.FileStoragePath(), fileName)
		request.Type = tmp.Type

		err = c.SaveUploadedFile(file, request.FileName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unable to save file"})
			return
		}

	case FileApiType:

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO: Check einbauen of Pfad existiert

	default:
		request.FileName = ""
	}

	if request.FileName != "" {

		task, err := h.service.CreateTask(
			c.Request.Context(),
			request.FileName,
			models.TaskType(request.Type),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, task)

	}
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := h.service.GetTask(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.service.GetAllTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var request struct {
		FileName string `json:"filename" binding:"required"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status := models.TaskStatus(request.Status)
	if status != models.StatusPending && status != models.StatusRunning &&
		status != models.StatusCompleted && status != models.StatusFailed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	task, err := h.service.UpdateTask(
		c.Request.Context(),
		uint(id),
		request.FileName,
		status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	if err := h.service.DeleteTask(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) DefaultReponse(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf(DefaultResponse, VERSION))
}
