package app

import (
	"net/http"

    "github.com/EgorikA4/golang-yandex-intensive-autumn-2024/config"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/handlers"
    "github.com/gorilla/mux"
)

type Application struct {
	orchestratorConfig *config.OrchestratorConfig
}

func New() *Application {
	return &Application{
		orchestratorConfig: config.GetOrchestratorConfig(),
	}
}

func (a *Application) RunServer() error {
    router := mux.NewRouter()
	router.HandleFunc("/api/v1/calculate/", handlers.CalcHandler)
    router.HandleFunc("/internal/task", handlers.TaskHandler)
    router.HandleFunc("/api/v1/expressions", handlers.GetExpressionListHandler)
    router.HandleFunc("/api/v1/expressions/{id}", handlers.GetExpressionHandler)
    http.Handle("/", router)
	return http.ListenAndServe(":" + a.orchestratorConfig.Port, nil)
}
