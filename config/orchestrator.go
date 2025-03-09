package config

import (
	"os"
	"strconv"
	"sync"
)

type OrchestratorConfig struct {
    Port string
}

var (
    orchestratorConfig *OrchestratorConfig
    orchestratorOnce sync.Once
)

func LoadOrchestratorConfig() error {
    var err error
    orchestratorOnce.Do(func() {
        orchestratorConfig = &OrchestratorConfig{}
        if _, err := strconv.Atoi(os.Getenv("ORCHESTRATOR_PORT")); err != nil {
            return
        }
        orchestratorConfig.Port = os.Getenv("ORCHESTRATOR_PORT")
    })
    return err
}
