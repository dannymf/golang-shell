package shell

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Pair struct {
	token, tokenType string
}

func MainLoop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("% ")
		// Read the keyboad input.
		input, err := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			continue
		}

		// Handle the execution of the input.
		lexed, err := Lexer(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			continue
		}

		parsed, err := Parser(lexed)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			continue
		}

		if err = ExecCommand(parsed); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			continue
		}
	}
}

// Split the input separate the command and the arguments.
// Return the command and the arguments.
func Lexer(input string) (*[]Pair, error) {
	var lexed *[]Pair = &[]Pair{}
	*lexed = make([]Pair, 0)

	a := []string{}

	sb := &strings.Builder{}
	quoted := false
	for _, r := range input {
		// change once able to handle redirects
		if !quoted && r == ' ' {
			a = append(a, sb.String())
			sb.Reset()
		} else if !IsValidToken(r) && !quoted {
			return nil, errors.New("cannot parse character " + string(r))
		} else if r == '"' || r == '\'' {
			quoted = !quoted
			// sb.WriteRune(r) // keep '"' otherwise comment this line
		} else if IsValidToken(r) || quoted {
			sb.WriteRune(r)
		} else {
			return nil, errors.New("cannot parse string " + string(input))
		}
	}
	if sb.Len() > 0 {
		a = append(a, sb.String())
	}

	for _, s := range a {
		if s == "<" {
			*lexed = append(*lexed, Pair{s, "STDIN-REDIRECT"})
		} else if s == ">" {
			*lexed = append(*lexed, Pair{s, "STDOUT-REDIRECT"})
		} else {
			*lexed = append(*lexed, Pair{s, "NORMAL"})
		}
	}
	// fmt.Println(lexed)
	return lexed, nil
}

// Parse the command and arguments into a list of tokens.
func Parser(lexed *[]Pair) (*[]Pair, error) {

	var parsed *[]Pair = &[]Pair{}
	*parsed = make([]Pair, 0)
	state := "START"
	var err error

	if err := ContainsMultipleRedirects(*lexed); err != nil {
		return nil, err
	}

	// Insert DFA here
	for _, pair := range *lexed {
		switch state {
		case "START":
			parsed, state, err = HandleStartState(parsed, pair, state)
			if err != nil {
				return nil, err
			}
		case "ARGUMENTS":
			parsed, state, err = HandleArgumentsState(parsed, pair, state)
			if err != nil {
				return nil, err
			}
		case "STDIN-REDIRECT":
			parsed, state, err = HandleStdinRedirectState(parsed, pair, state)
			if err != nil {
				return nil, err
			}
		case "STDOUT-REDIRECT":
			parsed, state, err = HandleStdoutRedirectState(parsed, pair, state)
			if err != nil {
				return nil, err
			}
		case "FILE-INPUT":
			parsed, state, err = HandleFileInputState(parsed, pair, state)
			if err != nil {
				return nil, err
			}
		}
	}

	if state == "STDOUT-REDIRECT" || state == "STDIN-REDIRECT" {
		return nil, errors.New("cannot end command with redirect")
	}

	return parsed, nil
}

func ExecCommand(parsed *[]Pair) error {

	if len(*parsed) == 0 {
		return nil
	}

	command := (*parsed)[0].token

	// Handle built in commands
	switch command {
	case "cd":

		if len(*parsed) < 2 {
			return os.Chdir(os.Getenv("HOME"))
		}
		return os.Chdir((*parsed)[1].token)

	case "exit", "quit":
		os.Exit(0)
	}

	// Get arguments from parsed
	args := make([]string, 0)
	for _, pair := range *parsed {
		if pair.tokenType == "ARGUMENT" {
			args = append(args, pair.token)
		}
	}

	var cmd *exec.Cmd

	if len(args) > 0 {
		cmd = exec.Command(command, args...)
	} else {
		cmd = exec.Command(command)
	}

	// Set the default output device-- DEFAULT
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Check if parsed has redirects
	if len(*parsed) > 1 {
		for idx, pair := range *parsed {
			if pair.tokenType == "STDIN-REDIRECT" {
				redirectFile, err := os.Open((*parsed)[idx+1].token)
				if err != nil {
					return err
				}
				cmd.Stdin = redirectFile
				defer redirectFile.Close()
			} else if pair.tokenType == "STDOUT-REDIRECT" {
				redirectFile, err := os.Create((*parsed)[idx+1].token)
				if err != nil {
					return err
				}
				cmd.Stdout = redirectFile
				defer redirectFile.Close()
			}
		}
	}
	// Run command with args
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
