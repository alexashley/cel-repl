package repl

import (
	"bufio"
	"fmt"
	"github.com/alexashley/cel-repl/repl/format"
	"github.com/alexashley/cel-repl/repl/history"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"log"
	"os"
	"strings"
)

const (
	quitFunctionId = "quit"

	inspectFunctionId      = "inspect"
	inspectFmtFunctionId   = "inspect_string"
	inspectAtFunctionId    = "inspect_int"
	inspectAtFmtFunctionId = "inspect_int_string"

	noHistory = types.String("No history to inspect")
)

var (
	defaultPrompt   = "> "
	globalFunctions = []*expr.Decl{
		decls.NewFunction(quitFunctionId, decls.NewOverload(quitFunctionId, []*expr.Type{}, decls.Null)),
		decls.NewFunction(inspectFunctionId, decls.NewOverload(inspectFunctionId, []*expr.Type{}, decls.String)),
		decls.NewFunction(inspectFunctionId, decls.NewOverload(inspectFmtFunctionId, []*expr.Type{decls.String}, decls.String)),
		decls.NewFunction(inspectFunctionId, decls.NewOverload(inspectAtFunctionId, []*expr.Type{decls.Int}, decls.String)),
		decls.NewFunction(inspectFunctionId, decls.NewOverload(inspectAtFmtFunctionId, []*expr.Type{decls.Int, decls.String}, decls.String)),
	}
)

type config struct {
	historySize    int
	check          bool
	macros         bool
	protoFilePaths []string
	prompt         string
}

type Repl struct {
	config    *config
	env       *cel.Env
	formatter *format.Formatter
	history   *history.EntryRingBuffer
	version   string
}

func NewRepl(version string) (*Repl, error) {
	c := &config{
		check:          true,
		macros:         true,
		prompt:         defaultPrompt,
		historySize:    100,
		protoFilePaths: []string{},
	}
	r := &Repl{
		config:    c,
		formatter: format.NewFormatter(),
		history:   history.NewEntryRingBuffer(c.historySize),
		version:   version,
	}

	opts := []cel.EnvOption{cel.Declarations(globalFunctions...)}
	if !c.macros {
		opts = append(opts, cel.ClearMacros())
	}

	env, err := cel.NewEnv(opts...)
	if err != nil {
		return nil, err
	}

	r.env = env

	return r, nil
}

func (r *Repl) quit(_ ...ref.Val) ref.Val {
	fmt.Println("So long, and thanks for all the fish!")
	os.Exit(0)

	return types.NullValue
}

func (r *Repl) inspect(_ ...ref.Val) ref.Val {
	entry := r.history.Get(r.history.Position() - 2)

	return r.inspectEntry(entry, format.Default)
}

func (r *Repl) inspectFmt(formatVal ref.Val) ref.Val {
	outputFormat := formatVal.Value().(string)
	entry := r.history.Get(r.history.Position() - 2)

	return r.inspectEntry(entry, outputFormat)
}

func (r *Repl) inspectAt(nthVal ref.Val) ref.Val {
	nth := nthVal.Value().(int64)

	return r.inspectEntry(r.history.Get(int(nth)-1), format.Default)
}

func (r *Repl) inspectAtFmt(nthVal ref.Val, formatVal ref.Val) ref.Val {
	nth := nthVal.Value().(int64)
	outputFormat := formatVal.Value().(string)

	return r.inspectEntry(r.history.Get(int(nth)-1), outputFormat)
}

func (r *Repl) inspectEntry(entry *history.Entry, format string) ref.Val {
	if entry == nil {
		return noHistory
	}

	if checkIssues(entry.Issues) {
		return types.String("that's not numberwang")
	}

	return types.String(r.formatter.FormatWith(entry.Ast.Expr(), format))
}

func (r *Repl) Init() *Repl {
	fmt.Printf("cel-repl %s started\n", strings.TrimSpace(r.version))
	fmt.Println("type quit() to exit")

	return r
}

func (r *Repl) Loop() {
	evalOptions := map[string]interface{}{}
	programOptions := cel.Functions(
		&functions.Overload{
			Operator: quitFunctionId,
			Function: r.quit,
		},
		&functions.Overload{
			Operator: inspectFunctionId,
			Function: r.inspect,
		},
		&functions.Overload{
			Operator: inspectFmtFunctionId,
			Unary:    r.inspectFmt,
		},
		&functions.Overload{
			Operator: inspectAtFunctionId,
			Unary:    r.inspectAt,
		},
		&functions.Overload{
			Operator: inspectAtFmtFunctionId,
			Binary:   r.inspectAtFmt,
		},
	)
	stdin := bufio.NewReader(os.Stdin)
	for {
		r.prompt()
		src, err := stdin.ReadString('\n')
		if err != nil {
			log.Fatal("Error reading stdin", err)
		}

		if len(strings.TrimSpace(src)) == 0 {
			continue
		}

		ast, issues := r.env.Compile(src)
		r.history.Insert(&history.Entry{
			Ast:    ast,
			Issues: issues,
			Raw:    src,
		})

		if checkIssues(issues) {
			r.displayIssues(ast, issues)
			continue
		}

		if r.config.check {
			ast, issues = r.env.Check(ast)
			if checkIssues(issues) {
				r.displayIssues(ast, issues)
				continue
			}
		}

		program, err := r.env.Program(ast, programOptions)
		if err != nil {
			r.prompt()
			fmt.Println(err)
			continue
		}

		result, _, err := program.Eval(evalOptions)
		if err != nil {
			r.prompt()
			fmt.Println(err)
			continue
		}

		r.prompt()
		fmt.Printf("%v\n", result.Value())
	}
}

func (r *Repl) prompt() {
	fmt.Printf(fmt.Sprintf("(%d)%s", r.history.Position()+1, r.config.prompt))
}

func newline() {
	fmt.Println()
}

func (r *Repl) getIssuesText(ast *cel.Ast, issues *cel.Issues) []string {
	var issueText []string
	for _, issue := range issues.Errors() {
		if ast != nil && ast.Source() != nil {
			issueText = append(issueText, issue.ToDisplayString(ast.Source()))
		} else {
			issueText = append(issueText, issue.Message)
		}
	}

	return issueText
}

func (r *Repl) displayIssues(ast *cel.Ast, issues *cel.Issues) {
	lines := r.getIssuesText(ast, issues)
	for _, line := range lines {
		r.prompt()
		fmt.Printf(strings.TrimSpace(line))
		newline()
	}
}

func checkIssues(issues *cel.Issues) bool {
	return issues != nil && len(issues.Errors()) > 0
}
