package utils

import (
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/config"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts"
)

func IsOperation(token string) bool {
    return token == "+" || token == "-" || token == "*" || token == "/"
}

func GetOperationTime(operation string) (int64, error) {
    cfg := config.GetAgentConfig()
    switch operation{
    case "+":
        return cfg.TimeAddMs, nil
    case "*":
        return cfg.TimeMultMs, nil
    case "-":
        return cfg.TimeSubMs, nil
    case "/":
        return cfg.TimeDivMs, nil
    }
    return 0, consts.ErrInvalidOperation
}

