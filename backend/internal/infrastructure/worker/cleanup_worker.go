package worker

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"

	"hufschlager.net/IngestListApiWrapper/internal/domain/tasks/interfaces"
	"hufschlager.net/IngestListApiWrapper/internal/domain/tasks/models"
)

type CleanUpWorker struct {
	repo            interfaces.TaskRepository
	semaphore       chan struct{}
	wg              sync.WaitGroup
	fileStoragePath string
}

func NewCleanUpWorker(
	repo interfaces.TaskRepository,
	maxWorkers int,
	storagePath string,
) *CleanUpWorker {
	return &CleanUpWorker{
		repo:            repo,
		semaphore:       make(chan struct{}, maxWorkers),
		fileStoragePath: storagePath,
	}
}

func (w *CleanUpWorker) ProcessTask(ctx context.Context, task *models.Task) {
	w.semaphore <- struct{}{}
	defer func() {
		<-w.semaphore
		w.wg.Done()
	}()

	fileName := task.FileName

	// check if the file was uploaded and stored locally
	if strings.Contains(fileName, w.fileStoragePath) {
		// check if file really exist
		_, err := os.Stat(fileName)
		if os.IsNotExist(err) {

			err = os.Remove(fileName)
			if err != nil {
				log.Printf("Error deleting file %s -> %v", fileName, err)
			}
		}
	}
	// Finally delete the task using the repo
	err := w.repo.Delete(ctx, task.ID)
	if err != nil {
		log.Printf("Error deleting task %s -> %v ", fileName, err)
	}
}

func (w *CleanUpWorker) StartTaskProcessing(ctx context.Context, task *models.Task) {
	w.wg.Add(1)
	go w.ProcessTask(ctx, task)
}

func (w *CleanUpWorker) Wait() {
	w.wg.Wait()
}
