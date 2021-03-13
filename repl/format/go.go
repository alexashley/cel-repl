package format

import "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

type goFormatter struct {}

func newGoFormatter() *goFormatter {
	return &goFormatter{}
}

func (g *goFormatter) Format(expression *expr.Expr) string {
	return expression.String()
}
