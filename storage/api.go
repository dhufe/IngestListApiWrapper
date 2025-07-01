package storage

import "dhufe/ingestlistapiwrapper/models"

type fileUploadResponse struct {
	DurationInMs int64  `json:"durationInMs"`
	FilePath     string `json:"filePath"`
}

type fileIndentify struct {
	DurationInMs int64 `json:"durationInMs"`
}

type fileIdentifyResponse struct {
	DurationInMs int64  `json:"durationInMs"`
	Result       string `json:"result"`
}

type fileIdentifyRequest struct {
	FilePath string `json:"filePath" binding:"required"`
}

type jobByIdResponse struct {
	Message string      `json:"message"`
	Data    models.Jobs `json:"data"`
}

type jobsResponse struct {
	Message string         `json:"message"`
	Data    *[]models.Jobs `json:"data"`
}
