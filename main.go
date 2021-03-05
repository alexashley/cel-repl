package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/google/cel-go/cel"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	//go:embed version
	replVersion   string
	defaultPrompt = "> "
	stripFromExit = regexp.MustCompile("\\.|:|[[:space:]]")
)

type config struct {
	check          bool
	macros         bool
	protoFilePaths []string
	prompt         string
}

type repl struct {
	config  *config
	env     *cel.Env
	history []string
}

func NewRepl() (*repl, error) {
	r := &repl{
		// TODO: allow overrides from file or flags
		config: &config{
			check:          true,
			macros:         true,
			prompt:         defaultPrompt,
			protoFilePaths: []string{},
		},
		history: []string{},
	}

	env, err := cel.NewEnv()
	if err != nil {
		return nil, err
	}

	r.env = env

	return r, nil
}

func (r *repl) init() {
	fmt.Printf("cel-repl %s started\n", strings.TrimSpace(replVersion))
}

func (r *repl) stop() {
	fmt.Println("So long, and thanks for all the fish!")
	os.Exit(0)
}

func (r *repl) handleCommand(src string) bool {
	src = strings.TrimSpace(src)

	src = stripFromExit.ReplaceAllString(src, "") // allow for :quit or .quit

	switch src {
	case "quit",
		"exit",
		"q":
		r.stop()
	case "inspect":
		if len(r.history) < 2 {
			fmt.Println("No history to inspect")
			return true
		}
		fmt.Println(len(r.history) - 1)
		fmt.Println(r.history)
		previous := r.history[len(r.history) - 2] // -2 because the `inspect` command is already in history
		ast, issues := r.env.Parse(previous)
		if checkIssues(issues) {
			r.displayIssues(ast, issues)
			return true
		}

		fmt.Println(ast.Expr())

		return true
	}

	return false
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

		r.history = append(r.history, src)

		if r.handleCommand(src) {
			continue
		}

		ast, issues := r.env.Compile(src)
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

		program, err := r.env.Program(ast)
		if err != nil {
			r.prompt()
			fmt.Print(err)
			continue
		}

		result, _, err := program.Eval(programOptions)
		if err != nil {
			r.prompt()
			fmt.Print(err)
			continue
		}

		r.prompt()
		fmt.Printf("%v\n", result.Value())
	}
}

func (r *repl) prompt() {
	fmt.Printf(r.config.prompt)
}

func newline() {
	fmt.Println()
}

func (r *repl) displayIssues(ast *cel.Ast, issues *cel.Issues) {
	if len(issues.Errors()) > 0 {
		for _, issue := range issues.Errors() {
			r.prompt()
			if ast != nil && ast.Source() != nil {
				fmt.Println(issue.ToDisplayString(ast.Source()))
			} else {
				fmt.Printf(issue.Message)
			}
			newline()
		}
	} else {
		r.prompt()
		fmt.Print(issues.Err())
		newline()
	}
}

func checkIssues(issues *cel.Issues) bool {
	return issues != nil && len(issues.Errors()) > 0
}

func main() {
	repl, err := NewRepl()
	if err != nil {
		log.Fatal("Failed to instantiate repl", err)
	}
	repl.init()
	repl.loop()
}
