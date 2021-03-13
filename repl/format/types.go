package format

import expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

type ExpressionFormatter interface {
	Format(expression *expr.Expr) string
}
