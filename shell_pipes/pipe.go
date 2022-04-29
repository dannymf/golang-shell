package shell_pipes

import (
	"os"
	"os/exec"
)

func Pipe() {
	// rd := bufio.NewReader(os.Stdin)
	file, _ := os.Create("file.txt")
	// file2, _ := os.Create("file2.txt")
	defer file.Close()
	// defer file2.Close()
	// r2 := bufio.NewReader(r2)

	cmd1 := exec.Command("cat", "dfa.txt")
	// cmd1.Stdin = r1
	cmd1.Stdout = file
	cmd1.Run()
	cmd1.Wait()

	cmd2 := exec.Command("wc")
	file2, _ := os.Open("file.txt")
	defer file2.Close()
	cmd2.Stdin = file2
	cmd2.Stdout = os.Stdout
	cmd2.Run()

	// rw.WriteTo(os.Stdout)
	// w1.Flush()

}

func ArbitraryPipe(parsed1 *[]Pair, parsed2 *[]Pair) error {

	interimFile, err := os.Create("interim.txt")
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
	interimFile2, _ := os.Open("interim.txt")
	cmd2.Stdin = interimFile2
	defer interimFile2.Close()
	cmd2.Stdout = os.Stdout
	// Delete interim file
	defer os.Remove("interim.txt")
	// Run commands
	if err := cmd1.Run(); err != nil {
		return err
	}
	if err := cmd2.Run(); err != nil {
		return err
	}
	return nil

}
