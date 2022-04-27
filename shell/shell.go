package shell

import (
	"bufio"
	"errors"
	"fmt"
	"go/parser"
	"os"
	"os/exec"
	"strings"
)

// func RunShell(cmd string) (string, error) {
// 	return "", nil
// }

type Pair struct {
	word, tokenType string
}

var ErrNoPath = errors.New("path required")

func MainLoop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("% ")
		// Read the keyboad input.
		input, err := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Handle the execution of the input.
		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func Lexer(input string) ([]Pair, error) {
	// Split the input separate the command and the arguments.
	// Return the command and the arguments.
	a := []string{}
	lexed := []Pair{}

	sb := &strings.Builder{}
	quoted := false
	for _, r := range input {
		// change once able to handle redirects
		if IsValidToken(r) || quoted {
			sb.WriteRune(r)
		} else if !IsValidToken(r) && !quoted {
			return nil, errors.New("cannot parse character " + string(r))
		} else if r == '"' || r == '\'' {
			quoted = !quoted
			// sb.WriteRune(r) // keep '"' otherwise comment this line
		} else if !quoted && r == ' ' {
			a = append(a, sb.String())
			sb.Reset()
		} else {
			return nil, errors.New("cannot parse string " + string(input))
		}
	}
	if sb.Len() > 0 {
		a = append(a, sb.String())
	}

	for _, s := range a {
		if s != "<" {
			lexed = append(lexed, Pair{s, "STDIN-REDIRECT"})
		} else if s != ">" {
			lexed = append(lexed, Pair{s, "STDOUT-REDIRECT"})
		} else {
			lexed = append(lexed, Pair{s, "NORMAL"})
		}
	}
	return lexed, nil
}

func Parser(lexed []Pair) {

}

func execInput(input string) error {
	// Remove the newline character.

	// Split the input separate the command and the arguments.
	command, args, err := parser(input)
	if err != nil {
		return err
	}

	// Check for built-in commands.
	switch command {
	case "cd":

		if len(args) == 0 || args[0] == "" {
			return os.Chdir(os.Getenv("HOME"))
		}
		return os.Chdir(args[1])

	case "exit":
		os.Exit(0)

	case "quit":
		os.Exit(0)
	}

	// Check if command is valid
	if err := IsValidCmd(command); err != nil {
		return err
	}

	//

	// // Check if command contains redirects in proper order
	// if indices, err := checkRedirectsOrder(args); err != nil {
	// 	return err
	// }

	// Check if command contains pipes
	// if containsPipe(args) {
	// 	return errors.New("invalid command")
	// }

	// Prepare the command to execute.
	cmd := exec.Command(command, args[0:]...)

	// Set the correct output device-- DEFAULT
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Check for redirects.
	if len(args) > 0 && (Contains(args, ">") || Contains(args, "<")) {
		// Check for stdout redirect.
		stdinContains, stdinIdx, stdinErr := ContainsRedirect(args, ">")
		if stdinErr != nil {
			return stdinErr
		}
		if stdinContains {
			// Open the file for writing.
			stdinErr = handleStdoutRedirect(args[stdinIdx+1], cmd)

			if stdinErr != nil {
				return stdinErr
			}
		}

		stdoutContains, stdoutIdx, stdoutErr := ContainsRedirect(args, "<")
		if stdoutErr != nil {
			return stdoutErr
		}
		if stdoutContains {
			// Open the file for writing.
			stdoutErr = handleStdinRedirect(args[stdoutIdx+1], cmd)

			if stdoutErr != nil {
				return stdoutErr
			}
		}

		// Execute the command and return the error.

	} else {
		return cmd.Run()
	}
	return nil
}

func syntacticAnalysis(command string, args []string) error {

	if command == "" || !IsOrdinaryString(command) {
		return errors.New("invalid command " + command)
	}
	return nil
}

func handleStdoutRedirect(filename string, cmd *exec.Cmd) error {
	// open the out file for writing
	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()
	return nil
}

func handleStdinRedirect(filename string, cmd *exec.Cmd) error {
	// open the out file for writing
	infile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer infile.Close()
	cmd.Stdin = infile

	err = cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()
	return nil
}
