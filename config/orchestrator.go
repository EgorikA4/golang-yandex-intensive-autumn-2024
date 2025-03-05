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
        orchestratorPort, ok := os.LookupEnv("APP_PORT")

        if !ok {
            orchestratorConfig.Port = "8000"
            return
        }

        if _, err := strconv.Atoi(orchestratorPort); err != nil {
            return
        }
        orchestratorConfig.Port = orchestratorPort
    })
    return err
}
