# COS316 Final Assignment: Recreating a Shell

## API
The code implements the following API

```go
type Pair struct {
	token, tokenType string
}

// Split the input into standard and special tokens
func Lexer(input string) (*[]Pair, error)

// Parse the command and arguments into a list of tokens.
// Return the command and the arguments.
func Parser(lexed *[]Pair) (*[]Pair, error)

// Execute the parsed command
func ExecCommand(parsed *[]Pair) error

```

The code includes a main client to run the terminal.

## Functionality:

The shell handles redirects from StdIn and to Stdout. The code does not support multiple redirect tokens.

There is also an additional main client (with more limited functionality) that implements the pipe syntax (synchronously).