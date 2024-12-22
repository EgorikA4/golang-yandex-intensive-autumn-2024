package calculation

import "errors"

var (
	ErrUnknownSymbol       = errors.New("unknown symbol")
	ErrDivisionByZero      = errors.New("division by zero")
	ErrEmptyExpression     = errors.New("empty expression")
	ErrMissingNum          = errors.New("missing number")
	ErrMissingOperation    = errors.New("missing operation")
	ErrMissingOpenBracket  = errors.New("there is no opening parenthesis corresponding to the closing one")
	ErrMissingCloseBracket = errors.New("there is no closing bracket corresponding to the opening one")
)