package shell

import (
	"os"
	"reflect"
	"testing"
)

var good1 string = "echo helloworld"
var good2 string = "cat helloworld.txt > output.txt"
var good3 string = "cat helloworld.txt wc -l"
var good4 string = "ls \"User/Desk top/Documents\""
var good5 string = "cat helloworld.txt > output.txt < input.txt"

var badCommand1 string = "cat non-existent-file.txt"

var goodCommand1 string = "touch new-file.txt"

var badLexer1 string = "cat helloworld.txt | output.txt < input.txt"
var badLexer2 string = "cat helloworld.txt & & | output.txt < input.txt"
var badLexer3 string = "cat helloworld.txt @ | output.txt < input.txt"
var badLexer4 string = "cat helloworld.txt ^ % | output.txt < input.txt"

var bad1 string = "cat helloworld.txt > output.txt > input.txt"
var bad2 string = "cat helloworld.txt < > output.txt < input.txt"
var bad3 string = "cat helloworld.txt output.txt > < input.txt"
var bad4 string = "echo < > helloworld.txt output.txt input.txt"
var bad5 string = "echo helloworld.txt output.txt input.txt >"

// var testString5 string =

var testStringsGood []string = []string{good1, good2, good3, good4, good5}
var testStringsBad []string = []string{bad1, bad2, bad3, bad4, bad5}
var testLexerBad []string = []string{badLexer1, badLexer2, badLexer3, badLexer4}

func TestLexer(testing *testing.T) {
	// Test the lexer.
	for _, testString := range testStringsGood {
		_, err := Lexer(testString)
		if err != nil {
			testing.Error(err)
		}
		// testing.Log(lexed)
		// for _, pair := range lexed {
		// 	fmt.Println(pair.token, pair.tokenType)
		// }
	}

}

func TestLexerBad(testing *testing.T) {
	// Test the lexer.
	for _, testString := range testLexerBad {
		_, err := Lexer(testString)

		if err == nil {
			testing.Error("error expected")
		}
	}
}

// }

func TestParser(testing *testing.T) {
	// Test the parser.
	for _, testString := range testStringsGood {
		lexed, err := Lexer(testString)
		if err != nil {
			testing.Error(err)
		}
		_, err = Parser(lexed)
		if err != nil {
			testing.Error(err)
		}
		// testing.Log(parsed)
		// fmt.Println(*parsed)
		// for _, pair := range *parsed {
		// 	fmt.Println(pair.token, pair.tokenType)
		// }
	}
}

func TestParserBad(testing *testing.T) {
	// Test the parser.
	for _, testString := range testStringsBad {
		lexed, err := Lexer(testString)
		if err != nil {
			testing.Error(err)
		}
		_, err = Parser(lexed)
		if err == nil {
			testing.Error("error expected")
		}
		// testing.Log(parsed)
		// fmt.Println(*parsed)
		// for _, pair := range *parsed {
		// 	fmt.Println(pair.token, pair.tokenType)
		// }
	}
}

func TestShellOutput(testing *testing.T) {
	// for _, testString := range testStringsGood {
	lexed, err := Lexer(good1)
	if err != nil {
		testing.Error(err)
	}
	parsed, err := Parser(lexed)
	if err != nil {
		testing.Error(err)
	}
	// testing.Log(parsed)
	// fmt.Println(*parsed)

	cmd, err := ReturnCommand(parsed)

	if err != nil {
		testing.Error(err)
	}

	output, err := cmd.Output()
	if err != nil {
		testing.Error(err)
	}
	correctOutput := []string{"helloworld"}

	if reflect.DeepEqual(output, correctOutput) {
		testing.Error("output is not correct")
	}
	// }
}

// Test cat command; file should not exist
func TestShellOutput2(testing *testing.T) {
	// for _, testString := range testStringsGood {
	lexed, err := Lexer(badCommand1)
	if err != nil {
		testing.Error(err)
	}
	parsed, err := Parser(lexed)
	if err != nil {
		testing.Error(err)
	}
	// testing.Log(parsed)
	// fmt.Println(*parsed)

	cmd, err := ReturnCommand(parsed)

	if err != nil {
		testing.Error(err)
	}

	err = cmd.Run()
	if err == nil {
		testing.Error("error expected")
	}

	if _, err = os.Stat("non-existent-file.txt"); err == nil {
		testing.Error("file should not exist")
	}
}

// Test touch command
func TestShellOutput3(testing *testing.T) {
	// for _, testString := range testStringsGood {
	lexed, err := Lexer(goodCommand1)
	if err != nil {
		testing.Error(err)
	}
	parsed, err := Parser(lexed)
	if err != nil {
		testing.Error(err)
	}
	// testing.Log(parsed)
	// fmt.Println(*parsed)

	cmd, err := ReturnCommand(parsed)

	if err != nil {
		testing.Error(err)
	}

	output, err := cmd.Output()
	if err != nil {
		testing.Error(err)
	}

	// No output should be returned
	correctOutput := []byte{}

	if !reflect.DeepEqual(output, correctOutput) {
		testing.Error("output is not correct")
	}

	if _, err = os.Stat("new-file.txt"); err != nil {
		testing.Error("file not successfully created")
	}

	os.Remove("new-file.txt")
}

// Test pwd command
func TestShellOutput4(testing *testing.T) {
	// for _, testString := range testStringsGood {
	lexed, err := Lexer("pwd")
	if err != nil {
		testing.Error(err)
	}
	parsed, err := Parser(lexed)
	if err != nil {
		testing.Error(err)
	}
	// testing.Log(parsed)
	// fmt.Println(*parsed)

	cmd, err := ReturnCommand(parsed)

	if err != nil {
		testing.Error(err)
	}

	output, err := cmd.Output()
	if err != nil {
		testing.Error(err)
	}

	// No output should be returned

	if err != nil {
		testing.Error(err)
	}

	// fmt.Println(string(output))

	// Change depending on where the test is run
	correctOutput := "/Users/Danny/Desktop/COS316/go_shell/shell\n"

	if string(output) != correctOutput {
		testing.Error("output is not correct")
	}
}
