package models

type AgentTask struct {
    Id string `json:"id"`
    Operation string `json:"operation"`
    OperationTime int64 `json:"operation_time"`
    Arg1 float64 `json:"arg1"`
    Arg2 float64 `json:"arg2"`
}

type AgentResult struct {
    Id string `json:"id"`
    Result float64 `json:"result"`
}

type Expression struct {
    Id string `json:"id"`
    Status string `json:"status"`
    Result float64 `json:"result"`
}

type ExpressionList struct {
    Expressions []Expression `json:"expressions"`
}

type Node struct {
    Id string
    ExpressionId string
    Operation string
    Arg1 string
    Arg2 string
}
