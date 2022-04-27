package shell

import "errors"

// throw error if contains multiple redirect tokens
func Contains(haystack []string, needle string) bool {
	for _, element := range haystack {
		if element == needle {
			return true
		}
	}
	return false
}

func ContainsRedirect(haystack []string, redirect string) (bool, int, error) {
	contained := false
	index := -1

	for idx, element := range haystack {
		if element == redirect && !contained {
			contained = true
			index = idx
		} else if element == redirect && contained {
			return false, -1, errors.New("multiple redirect tokens")
		}
	}
	if contained {
		return true, index, nil
	}
	return false, -1, nil
}

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

func IsOrdinaryToken(r rune) bool {
	return r != '<' && r != '>' && r != '|' && r != '&' && r != ';' && r != '(' && r != ')'
}

func IsValidToken(r rune) bool {
	return r != '|' && r != '&' && r != ';' && r != '(' && r != ')'
	// && r != '\'' && r != '"' && r != ' '
}

func IsOrdinaryString(s string) bool {
	for _, r := range s {
		if !IsOrdinaryToken(r) {
			return false
		}

	}
	return true
}

func IsValidCmd(s string) error {
	if s == "" || s == " " || s == "<" || s == ">" || s == "|" || s == "&" || s == ";" || s == "(" || s == ")" {
		return errors.New("invalid command")
	}
	return nil
}

// Check is > before <
// Returns true if > is before < and returns location of indices
func checkRedirectsOrder(args []string) ([]int, error) {

	bool1, int1, err1 := ContainsRedirect(args, "<")
	bool2, int2, err2 := ContainsRedirect(args, ">")

	if err1 != nil || err2 != nil {
		return nil, errors.New("invalid redirect")
	}

	if bool1 && bool2 {
		if int1 > int2 {
			return nil, errors.New("invalid redirect")
		}

	}
	return []int{int1, int2}, nil
}
