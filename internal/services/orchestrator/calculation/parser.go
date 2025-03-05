package calculation

import (
	"slices"
	"strconv"
	"strings"

	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/consts"
	"github.com/EgorikA4/golang-yandex-intensive-autumn-2024/internal/shared/utils"
)

func makePostfix(infixTokens []string) ([]string, error) {
	operations_priority := map[string]int{
		"+": 0,
		"-": 0,
		"*": 1,
		"/": 1,
	}
	brackets := []string{"(", ")"}
	var postfixTokens, stack []string
	prevType := "op"

	for _, token := range infixTokens {
		if utils.IsOperation(token) && prevType != "num" {
			return nil, consts.ErrMissingNum
		} else if utils.IsOperation(token) {
			for len(stack) > 0 && utils.IsOperation(stack[len(stack)-1]) && operations_priority[token] <= operations_priority[stack[len(stack)-1]] {
				postfixTokens = append(postfixTokens, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
			prevType = "op"
		} else if _, err := strconv.ParseFloat(token, 64); err == nil {
			if prevType == "num" {
				return nil, consts.ErrMissingOperation
			}
			postfixTokens = append(postfixTokens, token)
			prevType = "num"
		} else if slices.Contains(brackets, token) {
			if token == "(" {
				if prevType == "num" {
					return nil, consts.ErrMissingOperation
				}
				stack = append(stack, token)
				prevType = "bracket"
			} else {
				for len(stack) > 0 && stack[len(stack)-1] != "(" {
					postfixTokens = append(postfixTokens, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				}
				if len(stack) == 0 {
					return nil, consts.ErrMissingOpenBracket
				}
				stack = stack[:len(stack)-1]
			}
		} else {
			return nil, consts.ErrUnknownSymbol
		}
	}
	if prevType != "num" {
		return nil, consts.ErrMissingNum
	}
	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return nil, consts.ErrMissingCloseBracket
		}
		postfixTokens = append(postfixTokens, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return postfixTokens, nil
}

func Tokenize(expression string) ([]string, error) {
	if len(expression) == 0 {
		return nil, consts.ErrEmptyExpression
	}
	if expression[0] == '-' {
		expression = "0" + expression
	}
	expression = strings.ReplaceAll(expression, "(-", "(0-")
	symbols := [...]string{"+", "-", "*", "/", "(", ")"}
	for _, symbol := range symbols {
		expression = strings.ReplaceAll(expression, symbol, " "+symbol+" ")
	}

	tokens, err := makePostfix(strings.Fields(expression))
	if err != nil {
		return nil, err
	}
    return tokens, nil
}
