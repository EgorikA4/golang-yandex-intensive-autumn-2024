package config

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
)

var loadConfigOnce sync.Once

func LoadEnv() error {
	var err error
	loadConfigOnce.Do(func() {
		if err = godotenv.Load(); err != nil {
			return
		}

		err := LoadMemgraphConfig()
		if err != nil {
			err = fmt.Errorf("faild to load memgraph config: %v", err)
			return
		}

		err = LoadPostgresConfig()
		if err != nil {
			err = fmt.Errorf("faild to load postgres config: %v", err)
			return
		}

		err = LoadOrchestratorConfig()
		if err != nil {
			err = fmt.Errorf("faild to load orchestrator config: %v", err)
			return
		}

		err = LoadAgentConfig()
		if err != nil {
			err = fmt.Errorf("faild to load agent config: %v", err)
			return
		}
	})
	return err
}

func GetMemgraphConfig() *MemgraphConfig {
	return memgraphCfg
}

func GetPostgresConfig() *PostgresConfig {
	return postgresCfg
}

func GetOrchestratorConfig() *OrchestratorConfig {
	return orchestratorConfig
}

func GetAgentConfig() *AgentConfig {
	return agentConfig
}
