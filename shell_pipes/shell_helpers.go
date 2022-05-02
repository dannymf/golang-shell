package shell_pipes

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func IsValidToken(r rune) bool {
	return r != '&' && r != ';' && r != '(' && r != ')'
}

func ContainsMultiplePipes(lexed []Pair) error {
	pipeCount := false

	for _, pair := range lexed {
		if pair.token == "|" && !pipeCount {
			pipeCount = true
		} else if pair.token == "|" && pipeCount {
			return errors.New("multiple pipes")
		}
	}
	return nil
}

// Splits a Pair array into two slices before and after the pipe at index idx
func SplitOnPipe(parsed *[]Pair) (*[]Pair, *[]Pair) {
	var beforePipe []Pair = make([]Pair, 0)
	var afterPipe []Pair = make([]Pair, 0)

	var idx int

	// Find index of pipe in parsed
	for i, pair := range *parsed {
		if pair.token == "|" {
			idx = i
			break
		}
	}

	for i, pair := range *parsed {
		if i < idx {
			beforePipe = append(beforePipe, pair)
		} else if i > idx {
			afterPipe = append(afterPipe, pair)
		}
	}
	return &beforePipe, &afterPipe
}

func RunPipe(parsed1 *[]Pair, parsed2 *[]Pair) error {

	// Generate random number from 0 to 1000
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000)

	// Set filename to interim-file-<randomNumber>.txt
	filename := "interim-file-" + fmt.Sprint(randomNumber) + ".txt"
	interimFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer interimFile.Close()

	cmdStr1 := (*parsed1)[0].token
	cmdStr2 := (*parsed2)[0].token

	// Get arguments from parsed1
	args1 := make([]string, 0)
	for _, pair := range *parsed1 {
		if pair.tokenType == "ARGUMENT" {
			args1 = append(args1, pair.token)
		}
	}
	// Get arguments from parsed2
	args2 := make([]string, 0)
	for _, pair := range *parsed2 {
		if pair.tokenType == "ARGUMENT" {
			args2 = append(args2, pair.token)
		}
	}

	var cmd1 *exec.Cmd
	var cmd2 *exec.Cmd
	// Feed command and arguments into respective commands
	if len(args1) > 0 {
		cmd1 = exec.Command(cmdStr1, args1...)
	} else {
		cmd1 = exec.Command(cmdStr1)
	}
	if len(args2) > 0 {
		cmd2 = exec.Command(cmdStr2, args2...)
	} else {
		cmd2 = exec.Command(cmdStr2)
	}
	// Set input and output devices
	cmd1.Stdout = interimFile
	interimFile2, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer interimFile2.Close()

	cmd2.Stdin = interimFile2
	cmd2.Stdout = os.Stdout
	// Delete interim file
	defer os.Remove(filename)
	// Run commands
	if err := cmd1.Run(); err != nil {
		return err
	}
	if err := cmd2.Run(); err != nil {
		return err
	}
	return nil

}

// For testing purposes
// Does NOT handle cd or exit
func TestPipe(parsed1 *[]Pair, parsed2 *[]Pair) error {

	// Generate random number from 0 to 1000
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(1000)

	// Set filename to interim-file-<randomNumber>.txt
	filename := "interim-file-" + fmt.Sprint(randomNumber) + ".txt"
	interimFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer interimFile.Close()

	cmdStr1 := (*parsed1)[0].token
	cmdStr2 := (*parsed2)[0].token

	// Get arguments from parsed1
	args1 := make([]string, 0)
	for _, pair := range *parsed1 {
		if pair.tokenType == "ARGUMENT" {
			args1 = append(args1, pair.token)
		}
	}
	// Get arguments from parsed2
	args2 := make([]string, 0)
	for _, pair := range *parsed2 {
		if pair.tokenType == "ARGUMENT" {
			args2 = append(args2, pair.token)
		}
	}

	var cmd1 *exec.Cmd
	var cmd2 *exec.Cmd

	// Feed command and arguments into respective commands
	if len(args1) > 0 {
		cmd1 = exec.Command(cmdStr1, args1...)
	} else {
		cmd1 = exec.Command(cmdStr1)
	}
	if len(args2) > 0 {
		cmd2 = exec.Command(cmdStr2, args2...)
	} else {
		cmd2 = exec.Command(cmdStr2)
	}
	// Set input and output devices
	cmd1.Stdout = interimFile
	interimFile2, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer interimFile2.Close()

	cmd2.Stdin = interimFile2

	// Create temp file
	tempFile, err := os.Create("shell-temp.txt")
	defer tempFile.Close()
	cmd2.Stdout = tempFile

	// Delete interim file
	defer os.Remove(filename)
	// Run commands
	if err := cmd1.Run(); err != nil {
		return err
	}
	if err := cmd2.Run(); err != nil {
		return err
	}
	return nil
}
