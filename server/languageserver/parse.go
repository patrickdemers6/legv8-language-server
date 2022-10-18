package languageserver

import (
	"fmt"
	"math"
	"strings"

	lsp "go.lsp.dev/protocol"
)

func Parse(tokens *[]*[]*Token) *[]lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}

	// parse each line, reporting diagnostics when issues found
	for i, tokens := range *tokens {
		if len(*tokens) == 0 {
			continue
		}
		lineType := (*tokens)[0].Type

		// if there is a label token and colon token, this is a label. continue as no error found
		if lineType == LabelToken && len(*tokens) == 2 && (*tokens)[1].Type == ColonToken {
			continue
		}

		// since not a label, expect instruction
		if lineType != InstructionToken {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range: lsp.Range{
					Start: lsp.Position{Line: uint32(i), Character: 0},
					End:   lsp.Position{Line: uint32(i), Character: uint32(len((*tokens)[0].Value))},
				},
				Severity: lsp.DiagnosticSeverityError,
				Message:  "Expected an instruction keyword.",
				Source:   "compiler",
			})
			continue
		}

		instructionType := (*tokens)[0].InstructionType
		result := parse(tokens, i, expected[instructionType])
		if result != nil {
			diagnostics = append(diagnostics, *result)
		}
	}

	return &diagnostics
}

func parse(tokens *[]*Token, lineNumber int, expected *[]TokenType) *lsp.Diagnostic {
	for i, token := range *tokens {
		if i >= len(*expected) {
			return &lsp.Diagnostic{
				Range: lsp.Range{
					Start: lsp.Position{Line: uint32(lineNumber), Character: uint32(token.Start)},
					End:   lsp.Position{Line: uint32(lineNumber), Character: math.MaxUint32},
				},
				Severity: lsp.DiagnosticSeverityError,
				Message:  "Expected end of line.",
				Source:   "compiler",
			}
		}
		if token.Type != (*expected)[i] {
			return &lsp.Diagnostic{
				Range: lsp.Range{
					Start: lsp.Position{Line: uint32(lineNumber), Character: uint32(token.Start)},
					End:   lsp.Position{Line: uint32(lineNumber), Character: uint32(token.End)},
				},
				Severity: lsp.DiagnosticSeverityError,
				Message:  fmt.Sprintf("Expected a %s.", strings.ToLower((*expected)[i].String())),
				Source:   "compiler",
			}
		}
	}

	if len(*expected) != len(*tokens) {
		start := uint32((*tokens)[len(*tokens)-1].End)
		return &lsp.Diagnostic{
			Range: lsp.Range{
				Start: lsp.Position{Line: uint32(lineNumber), Character: start},
				End:   lsp.Position{Line: uint32(lineNumber), Character: start + 1},
			},
			Severity: lsp.DiagnosticSeverityError,
			Message:  fmt.Sprintf("Expected a %s.", strings.ToLower((*expected)[len(*tokens)].String())),
			Source:   "compiler",
		}
	}

	return nil
}

var expected map[InstructionType]*[]TokenType

func init() {
	expected = map[InstructionType](*[]TokenType){
		R:      &[]TokenType{InstructionToken, RegisterToken, CommaToken, RegisterToken, CommaToken, RegisterToken},
		I:      &[]TokenType{InstructionToken, RegisterToken, CommaToken, RegisterToken, CommaToken, NumberToken},
		IM:     &[]TokenType{InstructionToken, RegisterToken},
		D:      &[]TokenType{InstructionToken, RegisterToken, CommaToken, LeftBracketToken, RegisterToken, CommaToken, NumberToken, RightBracketToken},
		B:      &[]TokenType{InstructionToken, LabelToken},
		BR:     &[]TokenType{InstructionToken, RegisterToken},
		CB:     &[]TokenType{InstructionToken, RegisterToken, CommaToken, LabelToken},
		IGNORE: &[]TokenType{InstructionToken},
	}
}
