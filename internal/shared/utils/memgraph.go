package utils

import (
	"context"
	"fmt"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts/db_queries"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/models"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateNode(session neo4j.SessionWithContext, ctx context.Context, nodeId, operation, arg1, arg2, expressionId string) error {
	query := fmt.Sprintf(db_queries.CREATE_NODE, nodeId, operation, arg1, arg2, expressionId)
	if _, err := session.Run(ctx, query, nil); err != nil {
		return err
	}
	return nil
}

func CreateRightChildRel(session neo4j.SessionWithContext, ctx context.Context, childNodeId, parentNodeId string) error {
	query := fmt.Sprintf(db_queries.CREATE_RIGHT_REL, childNodeId, parentNodeId)
	if _, err := session.Run(ctx, query, nil); err != nil {
		return err
	}
	return nil
}

func CreateLeftChildRel(session neo4j.SessionWithContext, ctx context.Context, childNodeId, parentNodeId string) error {
	query := fmt.Sprintf(db_queries.CREATE_LEFT_REL, childNodeId, parentNodeId)
	if _, err := session.Run(ctx, query, nil); err != nil {
		return err
	}
	return nil
}

func GetIsRootAndExpressionId(session neo4j.SessionWithContext, ctx context.Context, nodeId string) (bool, string, error) {
	query := fmt.Sprintf(db_queries.IS_ROOT, nodeId)
	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return false, "", err
	}

	if !result.Next(ctx) {
		return false, "", consts.ErrNodeNotFound
	}

	record := result.Record()
	isRoot, _, err := neo4j.GetRecordValue[bool](record, "is_root")
	if err != nil {
		return false, "", err
	}

	expressionId, _, err := neo4j.GetRecordValue[string](record, "expression_id")
	if err != nil {
		return false, "", err
	}
	return isRoot, expressionId, nil
}

func DeleteNodeByExpressionId(session neo4j.SessionWithContext, ctx context.Context, expressionId string) error {
	query := fmt.Sprintf(db_queries.DELETE_NODE_BY_EXPRESSION_ID, expressionId)
	if _, err := session.Run(ctx, query, nil); err != nil {
		return err
	}
	return nil
}

func parseResult(result neo4j.ResultWithContext, ctx context.Context, nodeName string) (neo4j.Node, error) {
	record, err := result.Single(ctx)
	if err != nil {
		return neo4j.Node{}, err
	}

	node, isNil, err := neo4j.GetRecordValue[neo4j.Node](record, nodeName)
	if err != nil {
		return neo4j.Node{}, err
	}

	if isNil {
		return neo4j.Node{}, consts.ErrNodeNotFound
	}
	return node, nil
}

func parseNode(node neo4j.Node) (models.Node, error) {
	arg1, err := neo4j.GetProperty[string](node, "arg1")
	if err != nil {
		return models.Node{}, err
	}

	arg2, err := neo4j.GetProperty[string](node, "arg2")
	if err != nil {
		return models.Node{}, err
	}

	id, err := neo4j.GetProperty[string](node, "id")
	if err != nil {
		return models.Node{}, err
	}

	operation, err := neo4j.GetProperty[string](node, "operation")
	if err != nil {
		return models.Node{}, err
	}

	expressionId, err := neo4j.GetProperty[string](node, "expression_id")
	if err != nil {
		return models.Node{}, err
	}

	return models.Node{
		Id:           id,
		ExpressionId: expressionId,
		Operation:    operation,
		Arg1:         arg1,
		Arg2:         arg2,
	}, nil
}

func UpdateNode(session neo4j.SessionWithContext, ctx context.Context, nodeId string, nodeArg float64) (models.Node, error) {
	query := fmt.Sprintf(db_queries.UPDATE_NODE, nodeId, nodeArg, nodeArg)
	result, err := session.Run(ctx, query, nil)
	if err != nil {
		fmt.Println("error in the query!!!")
		return models.Node{}, err
	}

	node, err := parseResult(result, ctx, "parent")
	if err != nil {
		return models.Node{}, err
	}
	fmt.Println("parsed parent node", node)

	ans, err := parseNode(node)
	if err != nil {
		return models.Node{}, err
	}
	return ans, nil
}

func GetNodeById(session neo4j.SessionWithContext, ctx context.Context, nodeId string) (models.Node, error) {
	query := fmt.Sprintf(db_queries.GET_NODE_BY_ID, nodeId)
	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return models.Node{}, err
	}

	node, err := parseResult(result, ctx, "n")
	if err != nil {
		return models.Node{}, err
	}

	ans, err := parseNode(node)
	if err != nil {
		return models.Node{}, err
	}
	return ans, nil
}
