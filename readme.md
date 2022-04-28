# COS316 Final Assignment: Recreating a Shell

## API:
```go
// insert comment here
func execInput()

```

Handles redirect from StdIn, redirect from StdOut, Pipes, cd, exit, signals, 

First split on pipe, THEN check whether theres multiple redirects bc you can theoretically have piping and multiple redirects in each pipe

Do I need to use goroutines/channels?

First I will implement without pipes

Does NOT allow multiple redirects