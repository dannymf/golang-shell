package shell

import (
	"errors"
	"os"
	"os/exec"
)

func ContainsPipe(haystack []string) (bool, []int) {
	var indices []int = make([]int, 0)
	var found bool = false

	for idx, element := range haystack {
		if element == "|" {
			indices = append(indices, idx)
			found = true
		}
	}
	return found, indices
}

func IsValidToken(r rune) bool {
	return r != '|' && r != '&' && r != ';' && r != '(' && r != ')'
	// && r != '\'' && r != '"' && r != ' '
}

func ContainsMultipleRedirects(lexed []Pair) error {
	stdinRedirect := false
	stdoutRedirect := false

	for _, pair := range lexed {
		if pair.token == "<" && !stdinRedirect {
			stdinRedirect = true
		} else if pair.token == ">" && !stdoutRedirect {
			stdoutRedirect = true
		} else if (pair.token == "<" && stdinRedirect) || (pair.token == ">" && stdoutRedirect) {
			return errors.New("multiple redirects of same type")
		}
	}
	return nil
}

// For testing purposes
// Does NOT handle cd or exit
func ReturnCommand(parsed *[]Pair) (*exec.Cmd, error) {

	if len(*parsed) == 0 {
		return nil, nil
	}

	command := (*parsed)[0].token

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
	// cmd.Stdout = os.Stdout

	// Check if parsed has redirects
	if len(*parsed) > 1 {
		for idx, pair := range *parsed {
			if pair.tokenType == "STDIN-REDIRECT" {
				redirectFile, err := os.Open((*parsed)[idx+1].token)
				if err != nil {
					return nil, err
				}
				cmd.Stdin = redirectFile
				defer redirectFile.Close()
			} else if pair.tokenType == "STDOUT-REDIRECT" {
				redirectFile, err := os.Create((*parsed)[idx+1].token)
				if err != nil {
					return nil, err
				}
				cmd.Stdout = redirectFile
				defer redirectFile.Close()
			}
		}
	}
	// Run command with args
	return cmd, nil
}
