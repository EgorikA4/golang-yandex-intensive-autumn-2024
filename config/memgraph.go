package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/pkg/utils"
)

type MemgraphConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

var (
	memgraphCfg  *MemgraphConfig
	memgraphOnce sync.Once
)

func LoadMemgraphConfig() error {
	var err error
	memgraphOnce.Do(func() {
		requiredVars := []string{
			"MEMGRAPH_HOST",
			"MEMGRAPH_PORT",
			"MEMGRAPH_USERNAME",
			"MEMGRAPH_PASSWORD",
		}
		if err = utils.CheckEnvVars(requiredVars); err != nil {
			return
		}

		memgraphCfg = &MemgraphConfig{
			Host:     os.Getenv("MEMGRAPH_HOST"),
			Username: os.Getenv("MEMGRAPH_USERNAME"),
			Password: os.Getenv("MEMGRAPH_PASSWORD"),
		}

		memgraphPort := os.Getenv("MEMGRAPH_PORT")
		if _, err := strconv.Atoi(memgraphPort); err != nil {
			err = fmt.Errorf("MEMGRAPH_PORT should be integer")
			return
		}
		memgraphCfg.Port = memgraphPort
	})
	return err
}
