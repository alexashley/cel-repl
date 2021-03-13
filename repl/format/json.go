package format

import (
	"encoding/json"
	"google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type jsonFormatter struct {}

func newJsonFormatter() *jsonFormatter {
	return &jsonFormatter{}
}

func (j *jsonFormatter) Format(expression *expr.Expr) string {
	jsonBytes, _ := json.MarshalIndent(expression, "", "  ")

	return string(jsonBytes)
}
