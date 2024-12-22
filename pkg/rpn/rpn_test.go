package rpn_test

import (
	"testing"
	rpn "github.com/EgorikA4/golang-yandex-intensive-autumn-2024/pkg/rpn"
)

type Case struct {
	name string
	expression string
	expectedResult float64
	expectedError error
}

func TestCalc(t *testing.T) {
	cases := []Case{
		{
			name: "positive simple expression",
			expression: "2+2",
			expectedResult: 4,
			expectedError: nil,
		},
		{
			name: "positive expression with operations and priority",
			expression: "2+2*2",
			expectedResult: 6,
			expectedError: nil,
		},
		{
			name: "positive expression with float answer",
			expression: "1/2+2",
			expectedResult: 2.5,
			expectedError: nil,
		},
		{
			name: "positive expression with spaces",
			expression: "    2    +    2     ",
			expectedResult: 4,
			expectedError: nil,
		},
		{
			name: "positive expression with brackets priority",
			expression: "(2 + 2) * 2",
			expectedResult: 8,
			expectedError: nil,
		},
		{
			name: "positive expression with negative number next to a bracket",
			expression: "(-2) + (-2) * 2",
			expectedResult: -6,
			expectedError: nil,
		},
		{
			name: "positive expression with negative number at the beginning",
			expression: "-2 + 2 * 2",
			expectedResult: 2,
			expectedError: nil,
		},
		{
			name: "positie expression with float numbers",
			expression: "2.2 + 2.2",
			expectedResult: 4.4,
			expectedError: nil,
		},
		{
			name: "negative expression with ErrUnknownSymbol",
			expression: "2+2+?",
			expectedResult: 0,
			expectedError: rpn.ErrUnknownSymbol,
		},
		{
			name: "negative expression with ErrMissingOpenBracket",
			expression: "2 + 2 + 3)",
			expectedResult: 0,
			expectedError: rpn.ErrMissingOpenBracket,
		},
		{
			name: "negative expression with ErrMissingCloseBracket",
			expression: "(2 + 2 * 2",
			expectedResult: 0,
			expectedError: rpn.ErrMissingCloseBracket,
		},
		{
			name: "negative expression with ErrMissingNum",
			expression: "2 + + 2",
			expectedResult: 0,
			expectedError: rpn.ErrMissingNum,
		},
		{
			name: "negative expression with ErrMissingOperation",
			expression: "2 + 2 (2 * 2)",
			expectedResult: 0,
			expectedError: rpn.ErrMissingOperation,
		},
		{
			name: "negative expression with ErrMissingOperation",
			expression: "2 + 2 2",
			expectedResult: 0,
			expectedError: rpn.ErrMissingOperation,
		},
		{
			name: "negative empty expression",
			expression: "",
			expectedResult: 0,
			expectedError: rpn.ErrEmptyExpression,
		},
		{
			name: "negative expression with DivisionByZero",
			expression: "2 + 2 / 0",
			expectedResult: 0,
			expectedError: rpn.ErrDivisionByZero,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			gotNum, gotErr := rpn.Calc(testCase.expression)
			if !(gotNum == testCase.expectedResult && gotErr == testCase.expectedError) {
				t.Errorf(
					"Calc(%s) = %f, %v; expected %f, %v",
					testCase.expression,
					gotNum,
					gotErr,
					testCase.expectedResult,
					testCase.expectedError,
				)
			}
		})
	}
}