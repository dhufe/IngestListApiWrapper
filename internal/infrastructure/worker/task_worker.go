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
	repo       interfaces.TaskRepository
	outputRepo interfaces.TaskOutputRepository
	semaphore  chan struct{}
	wg         sync.WaitGroup
}

func NewTaskWorker(
	repo interfaces.TaskRepository,
	outputRepo interfaces.TaskOutputRepository,
	maxWorkers int,
) *TaskWorker {
	return &TaskWorker{
		repo:       repo,
		outputRepo: outputRepo,
		semaphore:  make(chan struct{}, maxWorkers),
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
	cmd := exec.CommandContext(ctx, task.Command)
	if task.Arguments != "" {
		cmd.Args = append(cmd.Args, task.Arguments)
	}

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

	// Ausgabe in separater Tabelle speichern
	outputRecord := &models.TaskOutput{
		TaskID: task.ID,
		Output: output,
	}
	if err := w.outputRepo.Create(ctx, outputRecord); err != nil {
		log.Printf("Error saving task output: %v", err)
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
