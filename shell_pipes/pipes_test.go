package shell_pipes

import (
	"os"
	"testing"
)

var testString1 string = `printf "Hello\n\nWorld\!\n" | wc`
var testOutput1 string = "       3       2      14\n"

var testString2 string = "echo COS316 | invalid"

func TestPipes(testing *testing.T) {

	lexed, err := Lexer(testString1)
	if err != nil {
		testing.Error(err)
	}
	parsed, err := Parser(lexed)
	if err != nil {
		testing.Error(err)
	}
	pair1, pair2 := SplitOnPipe(parsed)
	err = TestPipe(pair1, pair2)

	if err != nil {
		testing.Error(err)
	}

	tempFile, err := os.Open("shell-temp.txt")
	if err != nil {
		testing.Error(err)
	}

	var readBytes []byte = make([]byte, 25)
	tempFile.Read(readBytes)

	tempFile.Close()
	err = os.Remove("shell-temp.txt")

	if err != nil {
		testing.Error(err)
	}

	if string(readBytes) != testOutput1 {
		testing.Error("output is not correct")
	}

}

func TestPipes2(testing *testing.T) {

	lexed, err := Lexer(testString2)
	if err != nil {
		testing.Error(err)
	}
	parsed, err := Parser(lexed)
	if err != nil {
		testing.Error(err)
	}
	pair1, pair2 := SplitOnPipe(parsed)
	err = TestPipe(pair1, pair2)

	if err == nil {
		testing.Error("error expected")
	}
	os.Remove("shell-temp.txt")
}
