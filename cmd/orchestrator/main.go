package main

import (
	"log"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/config"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/app"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/storage"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/logger"
	"go.uber.org/zap"
)

func main() {
    logger, err := logger.InitLogger()
    if err != nil {
        log.Fatalf("can't initialize zap logger: %v", err)
        return
    }
    defer logger.Sync()
    logger = logger.With(zap.String("service", "orchestrator"))

    err = config.LoadEnv()
    if err != nil {
        logger.Error("faild to load config file", zap.Error(err))
        return
    }

    memgraphDB, err := storage.InitMemgraphDatabase()
    if err != nil {
        logger.Error("faild to connect to Memgraph database", zap.Error(err))
        return
    }
    defer memgraphDB.Close()

    postgresDB, err := storage.InitPostgresDatabase()
    if err != nil {
        logger.Error("faild to connect to Postgres database:", zap.Error(err))
        return
    }
    defer postgresDB.Close()

    logger.Info("orchestrator has been started")
	application := app.New()
	application.RunServer()
}
