package worker

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/tasks/interfaces"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/tasks/models"
	"github.com/dhufe/IngestListApiWrapper/pkg/utils"
)

type TaskWorker struct {
	repo      interfaces.TaskRepository
	semaphore chan struct{}
	wg        sync.WaitGroup
}

const (
	COMMAND = "third/bin/ingestlist"
)

func NewTaskWorker(
	repo interfaces.TaskRepository,
	maxWorkers int,
) *TaskWorker {
	return &TaskWorker{
		repo:      repo,
		semaphore: make(chan struct{}, maxWorkers),
	}
}

func (w *TaskWorker) ProcessTask(ctx context.Context, task *models.Task) {
	w.semaphore <- struct{}{}
	defer func() {
		<-w.semaphore
		w.wg.Done()
	}()

	task.Status = models.StatusRunning
	t := time.Now()
	task.StartedAt = &t
	if err := w.repo.Update(ctx, task); err != nil {
		log.Printf("Error updating task status to running: %v", err)
		return
	}

	args := []string{"identify", "-F", task.FileName}
	log.Printf("%s %s\n", COMMAND, args)
	// Befehl ausführen
	cmd := exec.CommandContext(ctx, COMMAND, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	convert := utils.NewXMLConverter(true, true)
	converted, err := convert.ToJSONFromBuffer(&stdout)
	task.Output = string(converted)
	task.Error = stderr.String()

	t = time.Now()
	task.CompletedAt = &t
	if err != nil {
		task.Status = models.StatusFailed
	} else {
		task.Status = models.StatusCompleted
	}
	// Task aktualisieren
	if err := w.repo.Update(ctx, task); err != nil {
		log.Printf("Error updating task after completion: %v", err)
	}
}

func (w *TaskWorker) StartTaskProcessing(ctx context.Context, task *models.Task) {
	w.wg.Add(1)
	go w.ProcessTask(ctx, task)
}

func (w *TaskWorker) Wait() {
	w.wg.Wait()
}
