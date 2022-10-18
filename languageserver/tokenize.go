package languageserver

import (
	"bufio"
	"math"
	"os"

	"go.lsp.dev/uri"
)

func TokenizeFile(file uri.URI) *[]*[]*Token {

	fileRead, err := os.Open(file.Filename())

	if err != nil {
		return nil
	}

	fileScanner := bufio.NewScanner(fileRead)

	fileScanner.Split(bufio.ScanLines)

	result := [](*[]*Token){}
	for fileScanner.Scan() {
		line := fileScanner.Text()
		tokens := TokenizeLine(line)
		result = append(result, tokens)
	}
	return &result
}

func TokenizeLine(line string) *[]*Token {

	var token *Token
	tokens := []*Token{}
	current := 0
	for current < len(line) {
		token, current = getNext(line, current)
		if token.Type == EOLToken {
			break
		}
		tokens = append(tokens, token)
	}

	return &tokens

}

func getNext(line string, current int) (*Token, int) {
	current = eatWhitespace(line, current)

	// if ends with whitespace
	if current >= len(line) {
		return &Token{
			Type:  EOLToken,
			Value: "",
		}, math.MaxInt32
	}

	result := handleSimpleTokens(line, current)
	if result != nil {
		return result, result.End
	}

	l := getNumber(line, current)
	if l > 0 {
		return &Token{
			Type:  NumberToken,
			Value: line[current : current+l],
			Start: current,
			End:   current + l,
		}, current + l
	}

	reg_len := isRegister(line, current)
	if reg_len > 0 {
		return &Token{
			Type:  RegisterToken,
			Value: line[current : current+reg_len],
			Start: current,
			End:   current + reg_len,
		}, current + reg_len
	}

	// check if the current position contains these instructions
	instructionsToCheck := []InstructionType{
		I, R, D, CB, IM, BR, IGNORE,
	}

	for _, h := range instructionsToCheck {
		result := h.Check(line, current)
		if result != nil {
			return result, result.End
		}
	}

	btype_len := isBType(line, current)
	if btype_len > 0 {
		return &Token{
			Type:            InstructionToken,
			InstructionType: B,
			Value:           line[current : current+btype_len],
			Start:           current,
			End:             current + btype_len,
		}, current + btype_len
	}

	ident := getIdentifier(line, current)
	if len(ident) > 0 {
		return &Token{
			Type:  LabelToken,
			Value: ident,
			Start: current,
			End:   current + len(ident),
		}, current + len(ident)
	}

	return &Token{
		Type:  UnknownToken,
		Value: string(line[current]),
		Start: current,
		End:   current + 1,
	}, current + 1

}

func handleSimpleTokens(line string, current int) *Token {
	switch line[current] {
	case ']':
		return &Token{
			Type:  RightBracketToken,
			Value: "]",
			Start: current,
			End:   current + 1,
		}
	case '[':
		return &Token{
			Type:  LeftBracketToken,
			Value: "[",
			Start: current,
			End:   current + 1,
		}
	case ',':
		return &Token{
			Type:  CommaToken,
			Value: ",",
			Start: current,
			End:   current + 1,
		}
	case ':':
		return &Token{
			Type:  ColonToken,
			Value: ":",
			Start: current,
			End:   current + 1,
		}
	case '/':
		if current+1 < len(line) && line[current+1] == '/' {
			return &Token{
				Type:  EOLToken,
				Value: "//",
				End:   math.MaxInt32,
			}
		}
	}
	return nil
}

func getIdentifier(line string, current int) string {
	end := current
	for end < len(line) && (isLetter(line[end]) || (end-current >= 1 && (line[end] == '_' || isNumber(line[end])))) {
		end++
	}
	return line[current:end]
}

func isRegister(line string, current int) int {
	// check registers X0-X9
	if current >= len(line)-1 {
		return 0
	}
	// X register
	if line[current] == 'X' {
		// X0-X9
		if line[current+1] >= '0' && line[current+1] <= '9' {
			// does reach end of line?
			if current >= len(line)-2 {
				return 2
			}
			// is next character 0-9?
			if line[current+2] >= '0' && line[current+2] <= '9' {
				return 3
			}
			return 2
		}
	}

	// handle SP, LR, FP
	if len(line)-1 > current {
		check := line[current : current+2]
		if check == "SP" || check == "FP" || check == "LR" {
			return 2
		}
	}

	// handel XZR
	if len(line)-2 > current && line[current:current+3] == "XZR" {
		return 3
	}

	return 0
}

func eatWhitespace(line string, current int) int {
	for current < len(line) && line[current] == ' ' {
		current++
	}
	return current
}

func isInstructionType(line string, current int, instructionType InstructionType) int {
	end := current
	for len(line) > end && isCapitalLetter(line[end]) {
		end++
	}
	if KeywordInstructionTypes[line[current:end]] == instructionType {
		return len(line[current:end])
	}
	return 0
}

func isBType(line string, current int) int {
	end := current

	for len(line) > end && (isCapitalLetter(line[end]) || (end-current == 1 && line[end] == '.')) {
		end++
	}

	if KeywordInstructionTypes[line[current:end]] == B {
		return len(line[current:end])
	}
	return 0
}

func isCapitalLetter(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func isLetter(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

func getNumber(line string, current int) int {
	end := current

	// at least one character on line after # and sequence starts with #
	if end >= len(line)-1 || line[end] != '#' {
		return 0
	}
	end++

	// iterate to end of number
	for end < len(line) && isNumber(line[end]) {
		end++
	}

	// just a pound sign, return 0 since not a number
	if end-current == 1 {
		return 0
	}

	// return number of characters in number
	return end - current
}

func isNumber(b byte) bool {
	return b >= '0' && b <= '9'
}
