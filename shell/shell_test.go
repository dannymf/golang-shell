package shell

import (
	"fmt"
	"testing"
)

var good1 string = "cat helloworld.txt"
var good2 string = "cat helloworld.txt > output.txt"
var good3 string = "cat helloworld.txt wc -l"
var good4 string = "ls \"User/Desk top/Documents\""
var good5 string = "cat helloworld.txt > output.txt < input.txt"

var badLexer string = "cat helloworld.txt | output.txt < input.txt"

var bad1 string = "cat helloworld.txt > output.txt > input.txt"
var bad2 string = "cat helloworld.txt < > output.txt < input.txt"
var bad3 string = "cat helloworld.txt output.txt > < input.txt"
var bad4 string = "echo < > helloworld.txt output.txt input.txt"
var bad5 string = "echo helloworld.txt output.txt input.txt >"

// var testString5 string =

var testStringsGood []string = []string{good1, good2, good3, good4, good5}
var testStringsBad []string = []string{bad1, bad2, bad3, bad4, bad5}

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
	_, err := Lexer(badLexer)
	if err == nil {
		testing.Error("error expected")
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
		parsed, err := Parser(lexed)
		if err != nil {
			testing.Error(err)
		}
		// testing.Log(parsed)
		fmt.Println(*parsed)
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
