package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/pkg/calculation"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/pkg/server"
)

type PositiveCase struct {
	name               string
	expression         string
	expectedResult   float64
	expectedStatusCode int
}

type NegativeCase struct {
	name string
	expression string
	expectedResult error
	expectedStatusCode int
}

func convertExpressionToBytesBuffer(expression string) *bytes.Buffer {
	bodyBytes, _ := json.Marshal(server.Request{
		Expression: expression,
	})
	return bytes.NewBuffer(bodyBytes)
}

func commonPositiveTestCalcHandler(t *testing.T, testCase PositiveCase) {
	w := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/calculate",
		convertExpressionToBytesBuffer(testCase.expression),
	)
	server.CalcHandler(w, req)

	result := w.Result()
	defer result.Body.Close()

	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	var response server.ResultResponse
	json.Unmarshal(bodyBytes, &response)
	if response.Result != testCase.expectedResult {
		t.Errorf(
			"CalcHandler(%s) = %v; expected %v",
			testCase.expression,
			response,
			testCase.expectedResult,
		)
	}

	if result.StatusCode != testCase.expectedStatusCode {
		t.Errorf(
			"CalcHandler(%s) -> StatusCode: %d; expected StatusCode: %d",
			testCase.expression,
			result.StatusCode,
			testCase.expectedStatusCode,
		)
	}
}

func commonNegativeTestCalcHandler(t *testing.T, testCase NegativeCase) {
	w := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/calculate",
		convertExpressionToBytesBuffer(testCase.expression),
	)
	server.CalcHandler(w, req)

	result := w.Result()
	defer result.Body.Close()

	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	var response server.ErrResponse
	json.Unmarshal(bodyBytes, &response)
	if response.Error != testCase.expectedResult.Error() {
		t.Errorf(
			"CalcHandler(%s) = %v; expected %v",
			testCase.expression,
			response,
			testCase.expectedResult,
		)
	}

	if result.StatusCode != testCase.expectedStatusCode {
		t.Errorf(
			"CalcHandler(%s) -> StatusCode: %d; expected StatusCode: %d",
			testCase.expression,
			result.StatusCode,
			testCase.expectedStatusCode,
		)
	}
}

func TestPositiveExpressionCalcHandler(t *testing.T) {
	cases := []PositiveCase{
		{
			name:           "simple expression",
			expression:     "2+2",
			expectedResult: 4,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "expression with operations and priority",
			expression:     "2+2*2",
			expectedResult: 6,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "expression with float answer",
			expression:     "1/2+2",
			expectedResult: 2.5,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "expression with spaces",
			expression:     "    2    +    2     ",
			expectedResult: 4,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "expression with brackets priority",
			expression:     "(2 + 2) * 2",
			expectedResult: 8,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "expression with negative number next to a bracket",
			expression:     "(-2) + (-2) * 2",
			expectedResult: -6,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "expression with negative number at the beginning",
			expression:     "-2 + 2 * 2",
			expectedResult: 2,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:           "expression with float numbers",
			expression:     "2.2 + 2.2",
			expectedResult: 4.4,
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			commonPositiveTestCalcHandler(t, testCase)
		})
	}
}

func TestNegativeExpressionCalcHandler(t *testing.T) {
	cases := []NegativeCase{
		{
			name:           "negative expression with ErrUnknownSymbol",
			expression:     "2+2+?",
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:  calculation.ErrUnknownSymbol,
		},
		{
			name:           "negative expression with ErrMissingOpenBracket",
			expression:     "2 + 2 + 3)",
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:  calculation.ErrMissingOpenBracket,
		},
		{
			name:           "negative expression with ErrMissingCloseBracket",
			expression:     "(2 + 2 * 2",
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:  calculation.ErrMissingCloseBracket,
		},
		{
			name:           "negative expression with ErrMissingNum",
			expression:     "2 + + 2",
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:  calculation.ErrMissingNum,
		},
		{
			name:           "negative expression with ErrMissingOperation",
			expression:     "2 + 2 (2 * 2)",
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:  calculation.ErrMissingOperation,
		},
		{
			name:           "negative expression with ErrMissingOperation",
			expression:     "2 + 2 2",
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:  calculation.ErrMissingOperation,
		},
		{
			name:           "negative empty expression",
			expression:     "",
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:  calculation.ErrEmptyExpression,
		},
		{
			name:           "negative expression with DivisionByZero",
			expression:     "2 + 2 / 0",
			expectedStatusCode: http.StatusBadRequest,
			expectedResult:  calculation.ErrDivisionByZero,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			commonNegativeTestCalcHandler(t, testCase)
		})
	}
}

func TestGetMethodCalcHandler(t *testing.T) {
	expression := "2 + 2"
	expectedResult := server.MethodNotAllowed
	expectedStatusCode := http.StatusMethodNotAllowed

	w := httptest.NewRecorder()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/calculate",
		convertExpressionToBytesBuffer(expression),
	)
	server.CalcHandler(w, req)

	result := w.Result()
	defer result.Body.Close()

	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	var response server.ErrResponse
	json.Unmarshal(bodyBytes, &response)
	if response.Error != expectedResult {
		t.Errorf(
			"CalcHandler(%s) = %v; expected %v",
			expression,
			response,
			expectedResult,
		)
	}

	if result.StatusCode != expectedStatusCode {
		t.Errorf(
			"CalcHandler(%s) -> StatusCode: %d; expected StatusCode: %d",
			expression,
			result.StatusCode,
			expectedStatusCode,
		)
	}
}