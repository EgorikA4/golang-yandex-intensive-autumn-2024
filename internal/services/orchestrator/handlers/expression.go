package handlers

import (
	"net/http"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/logger"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/services/orchestrator/storage"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func GetExpressionListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		return
	}

	logger := logger.GetLogger()

	pgInstance := storage.GetPostgresInstance()
	expressionList, err := utils.GetExpressionList(pgInstance.Conn, pgInstance.Ctx)
	if err != nil {
		logger.Error("failed to get the expressions list", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("the expressions list was successfully sent")
	utils.SendResponse(w, http.StatusOK, expressionList)
}

func GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		return
	}

	logger := logger.GetLogger()
	vars := mux.Vars(r)
	expressionId := vars["id"]

	if err := uuid.Validate(expressionId); err != nil {
		logger.Info("incorrect expression id", zap.String("expression_id", expressionId))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pgInstance := storage.GetPostgresInstance()
	expression, err := utils.GetExpression(pgInstance.Conn, pgInstance.Ctx, expressionId)
	if err == pgx.ErrNoRows {
		logger.Info("the expression was not found", zap.String("id", expressionId))
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error("faild to get the expression", zap.String("id", expressionId), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("the expression was successfully sent", zap.String("id", expressionId))
	utils.SendResponse(w, http.StatusOK, expression)
}
