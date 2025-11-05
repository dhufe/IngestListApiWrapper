package worker

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/tasks/interfaces"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/tasks/models"
)

type CleanUpWorker struct {
	repo      interfaces.TaskRepository
	semaphore chan struct{}
	wg        sync.WaitGroup
}

func NewCleanUpWorker(
	repo interfaces.TaskRepository,
	maxWorkers int,
) *CleanUpWorker {
	return &CleanUpWorker{
		repo:      repo,
		semaphore: make(chan struct{}, maxWorkers),
	}
}

func (w *CleanUpWorker) ProcessTask(ctx context.Context, task *models.Task) {
	w.semaphore <- struct{}{}
	defer func() {
		<-w.semaphore
		w.wg.Done()
	}()

	fileName := task.FileName

	// check if file exist
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {

		err = os.Remove(fileName)
		if err != nil {
			log.Printf("Error deleting file %s -> %v", fileName, err)
		}
		// Finally delete the task using the repo
		w.repo.Delete(ctx, task.ID)
	}
}

func (w *CleanUpWorker) StartTaskProcessing(ctx context.Context, task *models.Task) {
	w.wg.Add(1)
	go w.ProcessTask(ctx, task)
}

func (w *CleanUpWorker) Wait() {
	w.wg.Wait()
}
