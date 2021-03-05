## cel-repl

A bare-bones [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) for [cel](https://github.com/google/cel-spec), powered by [`cel-go`](https://github.com/google/cel-go).

```bash
$ go run main.go
cel-repl v0.1.0-local started
> 1 == 1
> true
> "foo".startsWith("bar")
> false
> 2 * 2
> 4
> [1, 2, 3].filter(n, n > 2)
> [3]
> quit
So long, and thanks for all the fish!
```

Feature progress:
- [x] repl
- [ ] inspect
    - [x] `inspect` shows the parsed expression of the last command
    - [ ] `inspect $n` does the same for the nth expression
    - [ ] `inspect $program` shows the parsed expression, but does not execute it
    - [ ] multiple formats
        - [ ] JSON
        - [x] default Go inspection
        - [ ] pretty printer that walks AST
- [ ] history
    - [x] store raw source
    - [ ] store previous parsed programs
    - [ ] configurable limit
    - [ ] up/down arrows to navigate
    - [ ] print history index in prompt
- [ ] load proto files and expose them to cel (probably only feasible if check is disabled)
- [ ] load config from a file or flags
- [ ] navigation shortcuts (ctrl-a/ctrl-e) for front
