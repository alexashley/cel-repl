package main

import (
	"bufio"
	_ "embed"
	"fmt"
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
	quitFunctionId      = "quit"
	inspectFunctionId   = "inspect"
	inspectAtFunctionId = "inspect_int"
	noHistory           = types.String("No history to inspect")
)

var (
	//go:embed version
	replVersion     string
	defaultPrompt   = "> "
	globalFunctions = []*expr.Decl{
		decls.NewFunction(quitFunctionId, decls.NewOverload(quitFunctionId, []*expr.Type{}, decls.Null)),
		decls.NewFunction(inspectFunctionId, decls.NewOverload(inspectFunctionId, []*expr.Type{}, decls.String)),
		decls.NewFunction(inspectFunctionId, decls.NewOverload(inspectAtFunctionId, []*expr.Type{decls.Int}, decls.String)),
	}
)

type config struct {
	history        int
	check          bool
	macros         bool
	protoFilePaths []string
	prompt         string
}

type repl struct {
	config  *config
	env     *cel.Env
	history *historyRingBuffer
}

func NewRepl() (*repl, error) {
	c := &config{
		check:          true,
		macros:         true,
		prompt:         defaultPrompt,
		history:        100,
		protoFilePaths: []string{},
	}
	r := &repl{
		config:  c,
		history: newHistoryRingBuffer(c.history),
	}

	env, err := cel.NewEnv(cel.Declarations(globalFunctions...))
	if err != nil {
		return nil, err
	}

	r.env = env

	return r, nil
}

func (r *repl) quit(_ ...ref.Val) ref.Val {
	fmt.Println("So long, and thanks for all the fish!")
	os.Exit(0)

	return types.NullValue
}

func (r *repl) inspect(_ ...ref.Val) ref.Val {
	entry := r.history.get(r.history.position() - 1)

	return r.inspectEntry(entry)
}

func (r *repl) inspectAt(nthVal ref.Val) ref.Val {
	nth := nthVal.Value().(int64)

	return r.inspectEntry(r.history.get(int(nth) - 1))
}

func (r *repl) inspectEntry(entry *historyEntry) ref.Val {
	if entry == nil {
		return noHistory
	}

	if checkIssues(entry.issues) {
		return types.String("that's not numberwang")
	}

	return types.String(entry.ast.Expr().String())
}

func (r *repl) init() {
	fmt.Printf("cel-repl %s started\n", strings.TrimSpace(replVersion))
	fmt.Println("type quit() to exit")
}

func (r *repl) loop() {
	programOptions := map[string]interface{}{}

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
		r.history.insert(&historyEntry{
			ast: ast,
			issues: issues,
			raw: src,
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

		program, err := r.env.Program(ast, cel.Functions(
			&functions.Overload{
				Operator: quitFunctionId,
				Function: r.quit,
			},
			&functions.Overload{
				Operator: inspectFunctionId,
				Function: r.inspect,
			},
			&functions.Overload{
				Operator: inspectAtFunctionId,
				Unary:    r.inspectAt,
			},
		))

		if err != nil {
			r.prompt()
			fmt.Println(err)
			continue
		}

		result, _, err := program.Eval(programOptions)
		if err != nil {
			r.prompt()
			fmt.Println(err)
			continue
		}

		r.prompt()
		fmt.Printf("%v\n", result.Value())
	}
}

func (r *repl) prompt() {
	fmt.Printf(fmt.Sprintf("(%d)%s", r.history.position()+1, r.config.prompt))
}

func newline() {
	fmt.Println()
}

func (r *repl) getIssuesText(ast *cel.Ast, issues *cel.Issues) []string {
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

func (r *repl) displayIssues(ast *cel.Ast, issues *cel.Issues) {
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
