package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/pkg/utils"
)

type AgentConfig struct {
	OrchestratorHost string
	OrchestratorPort string
	RequestsInterval int64
	NumWorkers       int
	TimeAddMs        int64
	TimeSubMs        int64
	TimeMultMs       int64
	TimeDivMs        int64
}

func parseTimeMs(varName string) (int64, error) {
	timeMs, err := strconv.ParseInt(os.Getenv(varName), 10, 64)
	if err != nil {
		return 0, err
	}

	if timeMs <= 0 {
		return 0, fmt.Errorf("%s=%d should be greater than zero", varName, timeMs)
	}
	return timeMs, nil
}

var (
	agentConfig *AgentConfig
	agentOnce   sync.Once
)

func LoadAgentConfig() error {
	var err error
	agentOnce.Do(func() {
		requiredVars := []string{
			"TIME_ADDITION_MS",
			"TIME_SUBTRACTION_MS",
			"TIME_MULTIPLICATIONS_MS",
			"TIME_DIVISIONS_MS",
			"COMPUTING_POWER",
			"ORCHESTRATOR_PORT",
			"ORCHESTRATOR_HOST",
		}

		if err = utils.CheckEnvVars(requiredVars); err != nil {
			return
		}

		agentConfig = &AgentConfig{}
		operationsTime := []string{"TIME_ADDITION_MS", "TIME_SUBTRACTION_MS", "TIME_MULTIPLICATIONS_MS", "TIME_DIVISIONS_MS"}
		for _, operationTime := range operationsTime {
			var timeMs int64
			timeMs, err = parseTimeMs(operationTime)
			if err != nil {
				return
			}

			switch operationTime {
			case "TIME_ADDITION_MS":
				agentConfig.TimeAddMs = timeMs
			case "TIME_SUBTRACTION_MS":
				agentConfig.TimeSubMs = timeMs
			case "TIME_MULTIPLICATIONS_MS":
				agentConfig.TimeMultMs = timeMs
			case "TIME_DIVISIONS_MS":
				agentConfig.TimeDivMs = timeMs
			default:
				err = fmt.Errorf("unknown field: %s in the AgentConfig struct", operationTime)
				return
			}
		}

		var numWorkers int
		numWorkers, err = strconv.Atoi(os.Getenv("COMPUTING_POWER"))
		if err != nil {
			return
		}
		if numWorkers <= 0 {
			err = fmt.Errorf("COMPUTING_POWER=%d should be greater than zero", numWorkers)
			return
		}
		agentConfig.NumWorkers = numWorkers

		agentConfig.OrchestratorHost = os.Getenv("ORCHESTRATOR_HOST")

		if _, err = strconv.Atoi(os.Getenv("ORCHESTRATOR_PORT")); err != nil {
			return
		}
		agentConfig.OrchestratorPort = os.Getenv("ORCHESTRATOR_PORT")

		if os.Getenv("REQUESTS_INTERVAL_MS") == "" {
			agentConfig.RequestsInterval = 1000
			return
		}

		var requestsInterval int64
		requestsInterval, err = parseTimeMs("REQUESTS_INTERVAL_MS")
		if err != nil {
			return
		}

		agentConfig.RequestsInterval = requestsInterval
	})
	return err
}
