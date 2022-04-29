package main

import (
	pipes "github.com/dannymf/golang-shell/shell_pipes"
)

func main() {
	// testString1 := "echo Joseph"
	// testString2 := "wc"

	// lexed1, _ := shell_pipes.Lexer(testString1)
	// parsed1, _ := shell_pipes.Parser(lexed1)

	// lexed2, _ := shell_pipes.Lexer(testString2)
	// parsed2, _ := shell_pipes.Parser(lexed2)
	// shell_pipes.ArbitraryPipe(parsed1, parsed2)
	pipes.MainLoop()
}
