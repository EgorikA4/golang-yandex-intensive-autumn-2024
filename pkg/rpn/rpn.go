package rpn

import (
	"errors"
	"slices"
	"strings"
	"strconv"
)

var (
	ErrUnknownSymbol = errors.New("unknown symbol")
	ErrDivisionByZero = errors.New("division by zero")
	ErrEmptyExpression = errors.New("empty expression")
	ErrMissingNum = errors.New("missing number")
	ErrMissingOperation = errors.New("missing operation")
	ErrMissingOpenBracket = errors.New("there is no closing bracket corresponding to the opening one")
	ErrMissingCloseBracket = errors.New("there is no opening parenthesis corresponding to the closing one")
)

func evalPostfix(postfixTokens []string) (float64, error) {
	var stack []float64
	for _, token := range postfixTokens {
		if parsedNum, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, parsedNum)
		} else {
			num2 := stack[len(stack) - 1]
			num1 := stack[len(stack) - 2]
			stack = stack[:len(stack) - 2]
			if token == "+" {
				stack = append(stack, num1 + num2)
			} else if token == "-" {
				stack = append(stack, num1 - num2)
			} else if token == "*" {
				stack = append(stack, num1 * num2)
			} else {
				if num2 == 0 {
					return 0, ErrDivisionByZero
				}
				stack = append(stack, num1 / num2)
			}
		}
	}
	return stack[0], nil
}

func makePostfix(infixTokens []string) ([]string, error) {
	operations := []string{"+", "-", "*", "/"}
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
		if slices.Contains(operations, token) && prevType != "num" {
			return nil, ErrMissingNum
		} else if slices.Contains(operations, token) {
			for len(stack) > 0 && slices.Contains(operations, stack[len(stack) - 1]) && operations_priority[token] <= operations_priority[stack[len(stack) - 1]] {
				postfixTokens = append(postfixTokens, stack[len(stack) - 1])
				stack = stack[:len(stack) - 1]
			}
			stack = append(stack, token)
			prevType = "op"
		} else if _, err := strconv.ParseFloat(token, 64); err == nil {
			if prevType == "num" {
				return nil, ErrMissingOperation
			}
			postfixTokens = append(postfixTokens, token)
			prevType = "num"
		} else if slices.Contains(brackets, token) {
			if token == "(" {
				if prevType == "num" {
					return nil, ErrMissingOperation
				}
				stack = append(stack, token)
				prevType = "bracket"
			} else {
				for len(stack) > 0 && stack[len(stack) - 1] != "(" {
					postfixTokens = append(postfixTokens, stack[len(stack) - 1])
					stack = stack[:len(stack) - 1]
				}
				if len(stack) == 0 {
					return nil, ErrMissingOpenBracket
				}
				stack = stack[:len(stack) - 1]
			}
		} else {
			return nil, ErrUnknownSymbol
		}
	}
	if prevType != "num" {
		return nil, ErrMissingNum
	}
	for len(stack) > 0 {
		if stack[len(stack) - 1] == "(" {
			return nil, ErrMissingCloseBracket
		}
		postfixTokens = append(postfixTokens, stack[len(stack) - 1])
		stack = stack[:len(stack) - 1]
	}
	return postfixTokens, nil
}

func Calc(expression string) (float64, error) {
	if len(expression) == 0 {
		return 0, ErrEmptyExpression
	}
	if expression[0] == '-' {
		expression = "0" + expression
	}
	expression = strings.ReplaceAll(expression, "(-", "(0-")
	symbols := [...]string{"+", "-", "*", "/", "(", ")"}
	for _, symbol := range symbols {
		expression = strings.ReplaceAll(expression, symbol, " " + symbol + " ")
	}
	tokens, err := makePostfix(strings.Fields(expression))
	if err != nil {
		return 0, err
	}
	ans, err := evalPostfix(tokens)
	if err != nil {
		return 0, err
	}
	return ans, nil
}
