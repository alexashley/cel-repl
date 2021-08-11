## cel-repl

test
A bare-bones [REPL](https://en.wikipedia.org/wiki/Read%E2%80%93eval%E2%80%93print_loop) for [cel](https://github.com/google/cel-spec), powered by [`cel-go`](https://github.com/google/cel-go).

```bash
$ make
cel-repl v0.1.0-local started
type quit() to exit
(1)> 1 == 1
(2)> true
(2)> "foo".startsWith("bar")
(3)> false
(3)> 2 * 2             
(4)> 4
(4)> [1, 2, 3].filter(n, n > 2)
(5)> [3]
(5)> inspect(1)
(6)> id:2  call_expr:{function:"_==_"  args:{id:1  const_expr:{int64_value:1}}  args:{id:3  const_expr:{int64_value:1}}}
(6)> quit()
So long, and thanks for all the fish!
```

Feature progress:
- [x] repl
- [ ] inspect
    - [x] `inspect()` shows the parsed expression of the last command
    - [x] `inspect(n)` does the same for the nth expression
    - [ ] multiple formats
        - [x] JSON
        - [x] default Go inspection
        - [ ] pretty printer that walks AST
- [ ] history
    - [x] store raw source
    - [x] store previous parsed programs
    - [x] configurable limit
    - [ ] up/down arrows to navigate
    - [x] print history index in prompt
- [x] expose REPL commands to cel instead of handling externally
- [ ] `compile` function: parse expression and show parsed output
    - [ ] `compile(src)`
    - [ ] `compile(src, fmt)` -- same format options as inspect
- [ ] `eval` function (compile/check/run)
- [ ] `ieval` function ("immediate eval" -- no check, compile and run)
- [ ] `help` function: print all available functions
- [ ] load proto files and expose them to cel
- [ ] load data from files -- eval with protos
- [ ] load config from a file/flags/environment
- [ ] navigation shortcuts (ctrl-a/ctrl-e)
- [ ] clear screen
- [ ] restart/reload session
- [ ] autocomplete
