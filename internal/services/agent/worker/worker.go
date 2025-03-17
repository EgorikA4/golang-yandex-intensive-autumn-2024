package worker

import (
	"time"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/models"
)

func performOperation(arg1, arg2 float64, operation string) float64 {
	if operation == "+" {
		return arg1 + arg2
	} else if operation == "-" {
		return arg1 - arg2
	} else if operation == "*" {
		return arg1 * arg2
	}
	return arg1 / arg2
}

func worker(tasks <-chan models.AgentTask, results chan<- models.AgentResult) {
	for task := range tasks {
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
		result := performOperation(task.Arg1, task.Arg2, task.Operation)
		results <- models.AgentResult{
			Id:     task.Id,
			Result: result,
		}
	}
}
