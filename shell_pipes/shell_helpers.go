package shell_pipes

import "errors"

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
	return r != '&' && r != ';' && r != '(' && r != ')'
	// && r != '\'' && r != '"' && r != ' '
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

func RunPipeCommand() {

}
