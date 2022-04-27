package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// shell.MainLoop()
	if err := exec.Command("<").Run(); err != nil {
		fmt.Println(err)
	}
}
