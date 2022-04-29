package shell

import "errors"

func HandleStartState(parsed *[]Pair, pair Pair, state string) (*[]Pair, string, error) {

	*parsed = append(*parsed, Pair{pair.token, "COMMAND"})
	// fmt.Println("PARSED: ", parsed)

	if pair.tokenType == "NORMAL" {
		state = "ARGUMENTS"
	} else if pair.tokenType == "STDIN-REDIRECT" {
		return nil, "", errors.New("cannot have stdin redirect in start state")
	} else if pair.tokenType == "STDOUT-REDIRECT" {
		return nil, "", errors.New("cannot have stdout redirect in start state")
	} else {
		return nil, "", errors.New("invalid next token")
	}
	return parsed, state, nil
}

func HandleArgumentsState(parsed *[]Pair, pair Pair, state string) (*[]Pair, string, error) {

	if pair.tokenType == "NORMAL" {
		*parsed = append(*parsed, Pair{pair.token, "ARGUMENT"})
		state = "ARGUMENTS" //redundant
	} else if pair.tokenType == "STDIN-REDIRECT" {
		*parsed = append(*parsed, Pair{pair.token, "STDIN-REDIRECT"})
		state = "STDIN-REDIRECT"
	} else if pair.tokenType == "STDOUT-REDIRECT" {
		*parsed = append(*parsed, Pair{pair.token, "STDOUT-REDIRECT"})
		state = "STDOUT-REDIRECT"
	} else {
		return nil, "", errors.New("invalid next token")
	}
	return parsed, state, nil

}

func HandleStdoutRedirectState(parsed *[]Pair, pair Pair, state string) (*[]Pair, string, error) {
	// We already checked for multiple redirects of same kind
	// fmt.Println("HANDLE STDOUT REDIRECT")
	if pair.tokenType == "STDIN-REDIRECT" {
		return nil, "", errors.New("cannot have stdin redirect in stdout redirect state")
	} else if pair.tokenType == "NORMAL" {
		*parsed = append(*parsed, Pair{pair.token, "FILE-INPUT"})
		state = "FILE-INPUT"
	}
	return parsed, state, nil
}

func HandleStdinRedirectState(parsed *[]Pair, pair Pair, state string) (*[]Pair, string, error) {
	// We already checked for multiple redirects of same kind
	// fmt.Println("HANDLE STDIN REDIRECT")
	if pair.tokenType == "STDOUT-REDIRECT" {
		return nil, "", errors.New("cannot have stdout redirect in stdin redirect state")
	} else if pair.tokenType == "NORMAL" {
		*parsed = append(*parsed, Pair{pair.token, "FILE-INPUT"})
		state = "FILE-INPUT"
	}
	return parsed, state, nil
}

func HandleFileInputState(parsed *[]Pair, pair Pair, state string) (*[]Pair, string, error) {
	if pair.tokenType == "NORMAL" {
		*parsed = append(*parsed, Pair{pair.token, "ARGUMENT"})
		state = "ARGUMENTS"
	} else if pair.tokenType == "STDIN-REDIRECT" {
		*parsed = append(*parsed, Pair{pair.token, "STDIN-REDIRECT"})
		state = "STDIN-REDIRECT"
	} else if pair.tokenType == "STDOUT-REDIRECT" {
		*parsed = append(*parsed, Pair{pair.token, "STDOUT-REDIRECT"})
		state = "STDOUT-REDIRECT"
	} else {
		return nil, "", errors.New("invalid next token")
	}
	return parsed, state, nil
}
