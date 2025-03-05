package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/logger"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/calculation"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/models"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/utils"

    "go.uber.org/zap"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    logger := logger.GetLogger()

    if r.Method == http.MethodGet {
        task, err := calculation.GetTask()
        if err == consts.ErrEmptyTasks {
            logger.Info("task not found")
            w.WriteHeader(http.StatusNotFound)
            return
        } else if err != nil {
            logger.Error("faild to get a task", zap.Error(err))
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        logger.Info(
            "the task was sent to the Agent",
            zap.String("id", task.Id),
            zap.String("operation", task.Operation),
            zap.Float64("arg1", task.Arg1),
            zap.Float64("arg2", task.Arg2),
        )
        utils.SendResponse(w, http.StatusOK, task)
    } else if r.Method == http.MethodPost {
        var agentResult models.AgentResult
        defer r.Body.Close()

        if err := json.NewDecoder(r.Body).Decode(&agentResult); err != nil {
            logger.Warn("incorrect response format from the Agent", zap.Error(err))
            w.WriteHeader(http.StatusUnprocessableEntity)
            return
        }

        if err := calculation.UpdateTree(agentResult); err == consts.ErrNodeNotFound {
            logger.Warn("the task was not found in the dependency tree", zap.String("id", agentResult.Id))
            w.WriteHeader(http.StatusNotFound)
            return
        } else if err != nil {
            logger.Error("failed to update the dependecy tree", zap.Error(err))
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        logger.Info("the task was successfully processed", zap.String("id", agentResult.Id))
        w.WriteHeader(http.StatusOK)
    } else {
        logger.Warn("method not allowed", zap.String("method", r.Method))
        w.WriteHeader(http.StatusMethodNotAllowed) 
    }
}
