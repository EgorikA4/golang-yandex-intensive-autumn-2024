package main

import (
	"fmt"
	"slices"
)

const ERROR_MSG = "Incorrect expression!"

func ConvertToPostfix(expression string, digits []rune) ([]rune, error) {
	priority := map[rune]int{
		'(': 0,
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}
	var operands, postfix []rune
	for _, item := range expression {
		if slices.Contains(digits, item) {
			postfix = append(postfix, item)
		} else if item == '(' {
			operands = append(operands, item)
		} else if item == ')' {
			fl := true
			for i := range len(operands) {
				index := len(operands) - i - 1
				postfix = append(postfix, operands[index])
				if operands[index] == '(' {
					operands = operands[:index]
					postfix = postfix[:len(postfix)-1]
					fl = false
					break
				}
			}
			if fl {
				return nil, fmt.Errorf("%s", ERROR_MSG)
			}
		} else {
			fl := true
			for i := range len(operands) {
				index := len(operands) - i - 1
				operand := operands[index]
				_, ok := priority[operand]

				if !ok {
					return nil, fmt.Errorf("%s", ERROR_MSG)
				}

				if priority[operand] >= priority[item] {
					postfix = append(postfix, operand)
				} else {
					operands = operands[:index+1]
					fl = false
					break
				}
			}
			if fl {
				operands = []rune{}
			}
			operands = append(operands, item)
		}
	}
	for i := range len(operands) {
		operand := operands[len(operands)-1-i]
		postfix = append(postfix, operand)
	}
	return postfix, nil
}

func Calc(expression string) (float64, error) {
	var digits []rune
	for digit := range 10 {
		digits = append(digits, rune(48+digit))
	}

	postfix, err := ConvertToPostfix(expression, digits)
	if err != nil {
		return 0, err
	}
	var stack []float64
	for _, item := range postfix {
		if slices.Contains(digits, item) {
			stack = append(stack, float64(item)-48)
		} else if item == ')' || item == '(' {
			return 0, fmt.Errorf("%s", ERROR_MSG)
		} else if len(stack) < 2 {
			return 0, fmt.Errorf("%s", ERROR_MSG)
		} else {
			stack_length := len(stack)
			num2 := stack[stack_length-1]
			num1 := stack[stack_length-2]
			if item == '/' && num2 == 0 {
				return 0, fmt.Errorf("%s", ERROR_MSG)
			}
			stack = stack[:stack_length-2]
			switch item {
			case '/':
				stack = append(stack, num1/num2)
			case '*':
				stack = append(stack, num1*num2)
			case '-':
				stack = append(stack, num1-num2)
			case '+':
				stack = append(stack, num1+num2)
			}
		}
	}
	if(len(stack) > 1 || len(stack) == 0){
		return 0, fmt.Errorf("%s", ERROR_MSG)
	}
	return stack[0], nil
}

func main() {
	res, err := Calc("")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
