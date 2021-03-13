package format

import (
	"fmt"
	"google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

const (
	Go      = "go"
	Json    = "json"
	Pretty  = "pretty"
	Default = Go
)

type Formatter struct {
	pretty ExpressionFormatter
	golang ExpressionFormatter
	json   ExpressionFormatter
}

func NewFormatter() *Formatter {
	return &Formatter{
		pretty: newPrettyFormatter(),
		golang: newGoFormatter(),
		json:   newJsonFormatter(),
	}
}

func (f *Formatter) Format(expression *expr.Expr) string {
	return f.FormatWith(expression, Default)
}

func (f *Formatter) FormatWith(expression *expr.Expr, format string) string {
	switch format {
	case Json:
		return f.json.Format(expression)
	case Go:
		return f.golang.Format(expression)
	case Pretty:
		return f.pretty.Format(expression)
	}

	return fmt.Sprintf("unrecognized output format %s", format)
}
