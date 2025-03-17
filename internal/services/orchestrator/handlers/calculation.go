package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/calculation"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/logger"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/storage"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ExpressionRequest struct {
	Expression string `json:"expression"`
}

type ExpressionIdResponse struct {
	Id string `json:"id"`
}

type ErrResponse struct {
	Error string `json:"error"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		utils.SendResponse(w, http.StatusMethodNotAllowed, ErrResponse{
			Error: consts.MethodNotAllowed,
		})
		return
	}

	var request ExpressionRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendResponse(w, http.StatusUnprocessableEntity, ErrResponse{
			Error: err.Error(),
		})
		return
	}

	logger := logger.GetLogger()

	expressionId := uuid.New().String()
	tokens, err := calculation.Tokenize(request.Expression)
	if err != nil {
		logger.Info("the expression contains invalid data", zap.String("expression", request.Expression))
		utils.SendResponse(w, http.StatusUnprocessableEntity, ErrResponse{
			Error: err.Error(),
		})
		return
	}

	pgInstance := storage.GetPostgresInstance()
	if err = utils.InsertExpression(pgInstance.Conn, pgInstance.Ctx, expressionId, consts.EXPRESSION_PROCESSING, 0.0); err != nil {
		logger.Error(consts.ErrQueryExecutionPostgres, zap.Error(err))
		utils.SendResponse(w, http.StatusInternalServerError, ErrResponse{
			Error: consts.InternalServerError,
		})
	}

	if err = calculation.BuildTree(tokens, expressionId); err != nil {
		logger.Error("failed to build a tree based on the expression", zap.Error(err))
		if err = utils.UpdateExpression(pgInstance.Conn, pgInstance.Ctx, expressionId, consts.EXPRESSION_FAILED, 0.0); err != nil {
			logger.Error(consts.ErrQueryExecutionPostgres, zap.Error(err))
		}
		utils.SendResponse(w, http.StatusInternalServerError, ErrResponse{
			Error: err.Error(),
		})
		return
	}

	utils.SendResponse(w, http.StatusOK, ExpressionIdResponse{
		Id: expressionId,
	})
}
