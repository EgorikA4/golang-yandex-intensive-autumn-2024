package listener

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/agent/logger"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/agent/worker"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/models"
	"go.uber.org/zap"
)

type TasksListener struct {
	interval   time.Duration
	client     *http.Client
	workerPool *worker.WorkerPool
}

func NewTasksListener(interval time.Duration, workerPool *worker.WorkerPool) *TasksListener {
	return &TasksListener{
		interval:   interval,
		client:     &http.Client{},
		workerPool: workerPool,
	}
}

func (tl *TasksListener) Listen(url string) {
	logger := logger.GetLogger()
	ticker := time.NewTicker(tl.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			resp, err := tl.client.Get(url)
			if err != nil || resp.StatusCode == http.StatusNotFound {
				logger.Info("task not found")
				continue
			}
			defer resp.Body.Close()

			var task models.AgentTask
			if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
				logger.Error("failed to decode the task from json", zap.Error(err))
				continue
			}

			logger.Info("the task has been sent to the pool", zap.String("id", task.Id))
			tl.workerPool.Tasks <- task
		}
	}
}

func (tl *TasksListener) Send(url string) {
	logger := logger.GetLogger()

	for {
		select {
		case agentResult := <-tl.workerPool.Results:
			body, err := json.Marshal(agentResult)
			if err != nil {
				logger.Error("failed to encode the task to json", zap.Error(err))
				continue
			}

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
			if err != nil {
				logger.Error("failed to send the task to the Orchestrator", zap.Error(err))
				continue
			}
			req.Header.Set("Content-Type", "application/json")
			tl.client.Do(req)
			logger.Info("the task was successfully sent to the Orchetrator", zap.String("id", agentResult.Id))
		}
	}
}
