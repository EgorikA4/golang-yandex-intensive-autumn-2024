package utils

import (
	"context"
	"fmt"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts/db_queries"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/models"
	"github.com/jackc/pgx/v5"
)

func InsertExpression(conn *pgx.Conn, ctx context.Context, expressionId, status string, result float64) error {
	query := fmt.Sprintf(db_queries.INSERT_EXPRESSION, expressionId, status, result)
	if _, err := conn.Exec(ctx, query); err != nil {
		return err
	}
	return nil
}

func UpdateExpression(conn *pgx.Conn, ctx context.Context, expressionId, status string, result float64) error {
	query := fmt.Sprintf(db_queries.UPDATE_EXPRESSION, status, result, expressionId)
	if _, err := conn.Exec(ctx, query); err != nil {
		return err
	}
	return nil
}

func GetExpression(conn *pgx.Conn, ctx context.Context, expressionId string) (models.Expression, error) {
	query := fmt.Sprintf(db_queries.GET_EXPRESSION_BY_ID, expressionId)

	var id, status string
	var result float64

	err := conn.QueryRow(ctx, query).Scan(&id, &status, &result)
	if err != nil {
		return models.Expression{}, err
	}
	return models.Expression{
		Id:     id,
		Status: status,
		Result: result,
	}, nil
}

func GetExpressionList(conn *pgx.Conn, ctx context.Context) (models.ExpressionList, error) {
	rows, err := conn.Query(ctx, db_queries.GET_EXPRESSIONS)
	if err != nil {
		return models.ExpressionList{}, err
	}
	defer rows.Close()

	expressionList := models.ExpressionList{}
	for rows.Next() {
		var expressionId, status string
		var result float64

		err = rows.Scan(&expressionId, &status, &result)
		if err != nil {
			return models.ExpressionList{}, err
		}

		expression := models.Expression{
			Id:     expressionId,
			Status: status,
			Result: result,
		}
		expressionList.Expressions = append(expressionList.Expressions, expression)
	}
	return expressionList, nil
}
