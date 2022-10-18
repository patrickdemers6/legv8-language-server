package languageserver

import (
	"math"
	"testing"

	lsp "go.lsp.dev/protocol"
)

func TestParse(t *testing.T) {
	inputs := []string{
		"ADDI X0 X1, #12",
		"SUBI X1, XZR, #9 uh oh",
		"ZZZ",
		"SUBI X0, X2, ",
		"LDUR SP, X2, #0]",
		"LDUR SP, [X2, #0",
		"B label_1 // comment",
		"B.EQ done",
		"BL link",
		"CBZ X1, top",
		"AND X12, X10, X1",
		"AND X12, X10, XZR",
		"AND X12, X10, SP",
		"LDUR SP, [X2, #0]",
	}

	expected_outs := []*lsp.Diagnostic{
		{
			Severity: lsp.DiagnosticSeverityError,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 8,
				},
				End: lsp.Position{
					Line:      0,
					Character: 10,
				},
			},
			Message: "Expected a comma.",
		},
		{
			Severity: lsp.DiagnosticSeverityError,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 17,
				},
				End: lsp.Position{
					Line:      0,
					Character: math.MaxUint32,
				},
			},
			Message: "Expected end of line.",
		},
		{
			Severity: lsp.DiagnosticSeverityError,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 0,
				},
				End: lsp.Position{
					Line:      0,
					Character: 3,
				},
			},
			Message: "Expected an instruction keyword.",
		},
		{
			Severity: lsp.DiagnosticSeverityError,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 12,
				},
				End: lsp.Position{
					Line:      0,
					Character: 13,
				},
			},
			Message: "Expected a immediate.",
		},
		{
			Severity: lsp.DiagnosticSeverityError,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 9,
				},
				End: lsp.Position{
					Line:      0,
					Character: 11,
				},
			},
			Message: "Expected a left bracket.",
		},
		{
			Severity: lsp.DiagnosticSeverityError,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 16,
				},
				End: lsp.Position{
					Line:      0,
					Character: 17,
				},
			},
			Message: "Expected a right bracket.",
		},
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}

	for i, in := range inputs {
		tokens := []*[]*Token{
			TokenizeLine(in),
		}
		out := Parse(&tokens)

		if len(*out) == 0 && expected_outs[i] != nil {
			t.Errorf("Issue not detected when parsing input: %s.", in)
			continue
		}

		if len(*out) > 0 && expected_outs[i] == nil {
			t.Errorf("Issue detected when no issue present when parsing input: %s. Out = %v", in, out)
			continue
		}
		if len(*out) == 0 && expected_outs[i] == nil {
			continue
		}

		if expected_outs[i].Message != (*out)[0].Message {
			t.Errorf("Expected message '%s'. Recieved '%s'. Input: %s", expected_outs[i].Message, (*out)[0].Message, in)
		}
		if expected_outs[i].Range.Start.Character != (*out)[0].Range.Start.Character {
			t.Errorf("Expected start character %d. Recieved %d. Input: %s", expected_outs[i].Range.Start.Character, (*out)[0].Range.Start.Character, in)
		}
		if expected_outs[i].Range.Start.Line != (*out)[0].Range.Start.Line {
			t.Errorf("Expected start line %d. Recieved %d. Input: %s", expected_outs[i].Range.Start.Line, (*out)[0].Range.Start.Line, in)
		}
		if expected_outs[i].Range.End.Character != (*out)[0].Range.End.Character {
			t.Errorf("Expected end character %d. Recieved %d. Input: %s", expected_outs[i].Range.End.Character, (*out)[0].Range.End.Character, in)
		}
		if expected_outs[i].Range.End.Line != (*out)[0].Range.End.Line {
			t.Errorf("Expected end line %d. Recieved %d. Input: %s", expected_outs[i].Range.End.Line, (*out)[0].Range.End.Line, in)
		}
	}
}
