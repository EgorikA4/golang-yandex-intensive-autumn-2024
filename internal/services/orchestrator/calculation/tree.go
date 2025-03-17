package calculation

import (
	"strconv"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/storage"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/models"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/utils"
	"github.com/google/uuid"
)

type pair struct {
	id    string
	token string
}

func BuildTree(tokens []string, expressionId string) error {
	memgraph := storage.GetDBInstance()
	session := memgraph.Session()
	defer session.Close(memgraph.Ctx)

	var stack []pair
	for _, token := range tokens {
		if utils.IsOperation(token) {
			right, left := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			parent := pair{
				id:    uuid.New().String(),
				token: token,
			}

			if err := utils.CreateNode(session, memgraph.Ctx, parent.id, parent.token, left.token, right.token, expressionId); err != nil {
				return err
			}

			if utils.IsOperation(right.token) {
				if err := utils.CreateRightChildRel(session, memgraph.Ctx, right.id, parent.id); err != nil {
					return err
				}
			}
			if utils.IsOperation(left.token) {
				if err := utils.CreateLeftChildRel(session, memgraph.Ctx, left.id, parent.id); err != nil {
					return err
				}
			}

			if !utils.IsOperation(left.token) && !utils.IsOperation(right.token) {
				arg1, err := strconv.ParseFloat(left.token, 64)
				if err != nil {
					return err
				}

				arg2, err := strconv.ParseFloat(right.token, 64)
				if err != nil {
					return err
				}

				operationTime, err := utils.GetOperationTime(parent.token)
				if err != nil {
					return err
				}

				tasks = append(tasks, models.AgentTask{
					Id:            parent.id,
					Operation:     parent.token,
					OperationTime: operationTime,
					Arg1:          arg1,
					Arg2:          arg2,
				})
			}

			stack = append(stack, parent)
		} else {
			stack = append(stack, pair{
				id:    "",
				token: token,
			})
		}
	}
	return nil
}

func UpdateTree(task models.AgentResult) error {
	memgraph := storage.GetDBInstance()
	session := memgraph.Session()
	defer session.Close(memgraph.Ctx)

	isRoot, expressionId, err := utils.GetIsRootAndExpressionId(session, memgraph.Ctx, task.Id)
	if err != nil {
		return err
	}

	if isRoot {
		pgInstance := storage.GetPostgresInstance()
		if err = utils.UpdateExpression(pgInstance.Conn, pgInstance.Ctx, expressionId, consts.EXPRESSION_COMPLETED, task.Result); err != nil {
			return err
		}

		if err = utils.DeleteNodeByExpressionId(session, memgraph.Ctx, expressionId); err != nil {
			return err
		}
		return nil
	}

	node, err := utils.UpdateNode(session, memgraph.Ctx, task.Id, task.Result)
	if err != nil {
		return err
	}

	if !utils.IsOperation(node.Arg1) && !utils.IsOperation(node.Arg2) {
		operationTime, err := utils.GetOperationTime(node.Operation)
		if err != nil {
			return err
		}

		arg1, err := strconv.ParseFloat(node.Arg1, 64)
		if err != nil {
			return err
		}

		arg2, err := strconv.ParseFloat(node.Arg2, 64)
		if err != nil {
			return err
		}

		tasks = append(tasks, models.AgentTask{
			Id:            node.Id,
			Operation:     node.Operation,
			OperationTime: operationTime,
			Arg1:          arg1,
			Arg2:          arg2,
		})
	}

	return nil
}
