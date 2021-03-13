package format

import (
	"fmt"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"strings"
)

type prettyFormatter struct {
	builder strings.Builder
}

func newPrettyFormatter() *prettyFormatter {
	return &prettyFormatter{}
}

func (p *prettyFormatter) Format(expression *expr.Expr) string {
	p.builder.Reset()

	p.pretty(expression, 0)

	return p.builder.String()
}


func (p *prettyFormatter) pretty(expression *expr.Expr, level int) {
	p.builder.WriteString(withIdent(level))
	switch expression.ExprKind.(type) {
	case *expr.Expr_IdentExpr:
		p.prettyIdent(expression, level)
	case *expr.Expr_SelectExpr:
		p.prettySelect(expression, level)
	case *expr.Expr_ConstExpr:
		p.prettyConst(expression, level)
	case *expr.Expr_CallExpr:
		p.prettyCall(expression, level)
	}
}

// foo
// (ident: foo)
func (p *prettyFormatter) prettyIdent(expression *expr.Expr, level int) {
	p.builder.WriteString("(ident: ")
	p.builder.WriteString(expression.GetIdentExpr().Name)
	p.builder.WriteString(")")
}

func (p *prettyFormatter) prettySelect(expression *expr.Expr, level int) {
	selectExpr := expression.GetSelectExpr()

	operand := selectExpr.Operand
	p.pretty(operand, level)

	field := selectExpr.Field

	p.builder.WriteString(".")
	p.builder.WriteString(field)
}

// 1
// (const<int64>: 1)
func (p *prettyFormatter) prettyConst(expression *expr.Expr, level int) {
	c := expression.GetConstExpr()
	var val interface{}
	p.builder.WriteString("(const<")
	switch c.ConstantKind.(type) {
	case *expr.Constant_BoolValue:
		p.builder.WriteString("bool")
		val = c.GetBoolValue()
	case *expr.Constant_DoubleValue:
		p.builder.WriteString("double")
		val = c.GetDoubleValue()
	case *expr.Constant_Int64Value:
		p.builder.WriteString("int64")
		val = c.GetInt64Value()
	case *expr.Constant_Uint64Value:
		p.builder.WriteString("uint64")
		val = c.GetUint64Value()
	case *expr.Constant_StringValue:
		p.builder.WriteString("string")
		val = c.GetStringValue()
	default:
		val = c.GetNullValue()
		p.builder.WriteString("unknown")
	}

	p.builder.WriteString(">: ")
	p.builder.WriteString(fmt.Sprintf("%v", val))
	p.builder.WriteString(")")
}

// 1 + 1
// (call: +)
// 	(const<int64>: 1)
// 	(const<int64>: 1)
func (p *prettyFormatter) prettyCall(expression *expr.Expr, level int) {
	callExpr := expression.GetCallExpr()

	p.builder.WriteString("(call: ")
	p.builder.WriteString(callExpr.Function)
	p.builder.WriteString(")")

	// TODO: handle target
	for i := range callExpr.Args {
		p.pretty(callExpr.Args[i], level+1)
	}
}

func withIdent(level int) string {
	if level == 0 {
		return ""
	}

	indent := "\n"

	for i := 0; i < level; i++ {
		indent += "\t"
	}

	return indent
}
