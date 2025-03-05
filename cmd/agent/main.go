package main

import (
	"fmt"
	"log"
	"time"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/config"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/agent/listener"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/agent/logger"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/agent/worker"
	"go.uber.org/zap"
)

func main() {
    logger, err := logger.InitLogger()
    if err != nil {
        log.Fatalf("can't initialize zap logger: %v", err)
        return
    }
    defer logger.Sync()
    logger = logger.With(zap.String("service", "agent")) 

    err = config.LoadEnv()
    if err != nil {
        logger.Error("faild to load config file", zap.Error(err))
        return
    }

    orchestratorCfg := config.GetOrchestratorConfig()
    agentCfg := config.GetAgentConfig()
    url := fmt.Sprintf("http://localhost:%s/internal/task", orchestratorCfg.Port)

    workerPool := worker.NewWorkerPool(agentCfg.NumWorkers)
    taskListener := listener.NewTasksListener(time.Second, workerPool)

    logger.Info("agent has been started")
    go taskListener.Send(url)
    taskListener.Listen(url)
}
