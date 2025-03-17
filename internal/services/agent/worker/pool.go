package worker

import (
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/agent/logger"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/models"
	"go.uber.org/zap"
)

type WorkerPool struct {
	Tasks   chan models.AgentTask
	Results chan models.AgentResult
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	logger := logger.GetLogger()
	pool := &WorkerPool{
		Tasks:   make(chan models.AgentTask, numWorkers),
		Results: make(chan models.AgentResult, numWorkers),
	}

	for i := 1; i <= numWorkers; i++ {
		go worker(pool.Tasks, pool.Results)
		logger.Info("worker has been started", zap.Int("id", i))
	}
	return pool
}
