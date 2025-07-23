package worker

import (
	"bytes"
	"context"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/dhufe/IngestListApiWrapper/internal/domain/interfaces"
	"github.com/dhufe/IngestListApiWrapper/internal/domain/models"
)

type TaskWorker struct {
	repo      interfaces.TaskRepository
	semaphore chan struct{}
	wg        sync.WaitGroup
}

const (
	COMMAND = ""
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

	// Befehl ausführen
	cmd := exec.CommandContext(ctx, COMMAND)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()
	stderrOutput := stderr.String()

	t = time.Now()
	task.CompletedAt = &t
	if err != nil {
		task.Status = models.StatusFailed
		task.Error = stderrOutput
	} else {
		task.Status = models.StatusCompleted
	}

	// Ausgabe speichern
	task.Output = output

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
