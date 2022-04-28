package shell

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// func RunShell(cmd string) (string, error) {
// 	return "", nil
// }

type Pair struct {
	token, tokenType string
}

var lexed *[]Pair = &[]Pair{}
var parsed *[]Pair = &[]Pair{}
var state string = "START"

// var ErrNoPath = errors.New("path required")

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

		if err = execCommand(parsed); err != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", err)
			continue
		}
	}
}

func Lexer(input string) (*[]Pair, error) {
	// Split the input separate the command and the arguments.
	// Return the command and the arguments.
	lexed = &[]Pair{}
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

	// for _, s := range a {
	// 	fmt.Println(s)
	// }

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

func Parser(lexed *[]Pair) (parsed2 *[]Pair, err error) {
	// Check for incorrect order of redirects
	*parsed = make([]Pair, 0)
	state = "START"

	if err := ContainsMultipleRedirects(*lexed); err != nil {
		return nil, err
	}

	// fmt.Println("LEXED", *lexed)
	// Insert DFA here
	for _, pair := range *lexed {
		switch state {
		case "START":
			// fmt.Println("START")
			if err := HandleStartState(pair); err != nil {
				return nil, err
			}
		case "ARGUMENTS":
			// fmt.Println("ARGUMENTS")
			if err := HandleArgumentsState(pair); err != nil {
				return nil, err
			}
		case "STDIN-REDIRECT":
			// fmt.Println("STDIN-REDIRECT")
			if err := HandleStdinRedirectState(pair); err != nil {
				return nil, err
			}
		case "STDOUT-REDIRECT":
			// fmt.Println("STDOUT-REDIRECT")
			if err := HandleStdoutRedirectState(pair); err != nil {
				return nil, err
			}
		case "FILE-INPUT":
			// fmt.Println("FILE-INPUT")
			if err := HandleFileInputState(pair); err != nil {
				return nil, err
			}
		}
		// fmt.Println("PARSED LOOP: ", parsed)
	}

	if state == "STDOUT-REDIRECT" || state == "STDIN-REDIRECT" {
		return nil, errors.New("cannot end command with redirect")
	}

	// fmt.Println("PARSED POST: ", parsed)
	return parsed, nil
}

func execCommand(parsed *[]Pair) error {
	// Get command string

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

	case "exit":
		os.Exit(0)

	case "quit":
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

func HandleStartState(pair Pair) error {

	*parsed = append(*parsed, Pair{pair.token, "COMMAND"})
	// fmt.Println("PARSED: ", parsed)

	if pair.tokenType == "NORMAL" {
		state = "ARGUMENTS"
	} else if pair.tokenType == "STDIN-REDIRECT" {
		return errors.New("cannot have stdin redirect in start state")
	} else if pair.tokenType == "STDOUT-REDIRECT" {
		return errors.New("cannot have stdout redirect in start state")
	} else {
		return errors.New("invalid next token")
	}
	return nil
}

func HandleArgumentsState(pair Pair) error {

	if pair.tokenType == "NORMAL" {
		*parsed = append(*parsed, Pair{pair.token, "ARGUMENT"})
		state = "ARGUMENTS" //redundant
	} else if pair.tokenType == "STDIN-REDIRECT" {
		*parsed = append(*parsed, Pair{pair.token, "STDIN-REDIRECT"})
		state = "STDIN-REDIRECT"
	} else if pair.tokenType == "STDOUT-REDIRECT" {
		*parsed = append(*parsed, Pair{pair.token, "STDOUT-REDIRECT"})
		state = "STDOUT-REDIRECT"
	} else {
		return errors.New("invalid next token")
	}
	// fmt.Println("PARSED: ", parsed)
	return nil

}

func HandleStdoutRedirectState(pair Pair) error {
	// We already checked for multiple redirects of same kind
	// fmt.Println("HANDLE STDOUT REDIRECT")
	if pair.tokenType == "STDIN-REDIRECT" {
		return errors.New("cannot have stdin redirect in stdout redirect state")
	} else if pair.tokenType == "NORMAL" {
		*parsed = append(*parsed, Pair{pair.token, "FILE-INPUT"})
		state = "FILE-INPUT"
	}
	return nil
}

func HandleStdinRedirectState(pair Pair) error {
	// We already checked for multiple redirects of same kind
	// fmt.Println("HANDLE STDIN REDIRECT")
	if pair.tokenType == "STDOUT-REDIRECT" {
		return errors.New("cannot have stdout redirect in stdin redirect state")
	} else if pair.tokenType == "NORMAL" {
		*parsed = append(*parsed, Pair{pair.token, "FILE-INPUT"})
		state = "FILE-INPUT"
	}
	return nil
}

func HandleFileInputState(pair Pair) error {
	if pair.tokenType == "NORMAL" {
		*parsed = append(*parsed, Pair{pair.token, "ARGUMENT"})
		state = "ARGUMENTS"
	} else if pair.tokenType == "STDIN-REDIRECT" {
		*parsed = append(*parsed, Pair{pair.token, "STDIN-REDIRECT"})
		state = "STDIN-REDIRECT"
	} else if pair.tokenType == "STDOUT-REDIRECT" {
		*parsed = append(*parsed, Pair{pair.token, "STDOUT-REDIRECT"})
		state = "STDOUT-REDIRECT"
	} else {
		return errors.New("invalid next token")
	}
	return nil
}
