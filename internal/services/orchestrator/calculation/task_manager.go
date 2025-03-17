package calculation

import (
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/storage"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/models"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/utils"
)

var tasks []models.AgentTask

func GetTask() (models.AgentTask, error) {
	memgraph := storage.GetDBInstance()
	session := memgraph.Session()
	defer session.Close(memgraph.Ctx)

	pgInstance := storage.GetPostgresInstance()

	for len(tasks) > 0 {
		task := tasks[0]
		tasks = tasks[1:]
		if task.Operation == "/" && task.Arg2 == 0 {
			node, err := utils.GetNodeById(session, memgraph.Ctx, task.Id)
			if err != nil {
				return models.AgentTask{}, err
			}

			if err := utils.DeleteNodeByExpressionId(session, memgraph.Ctx, node.ExpressionId); err != nil {
				return models.AgentTask{}, err
			}

			if err := utils.UpdateExpression(pgInstance.Conn, pgInstance.Ctx, node.ExpressionId, consts.EXPRESSION_FAILED, 0.0); err != nil {
				return models.AgentTask{}, err
			}
			continue
		}
		return task, nil
	}
	return models.AgentTask{}, consts.ErrEmptyTasks
}
