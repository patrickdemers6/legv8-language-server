package languageserver

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	inputs := []string{
		"ADDI X0, X1, #12",
		"SUBI X1, XZR, #9",
		"LDUR SP, [X2, #0]",
		"B label_1",
		"B.EQ done",
		"BL link",
		"CBZ X1, top",
		"AND X12, X10, X1",
		"AND X12, X10, XZR",
		"AND X12, X10, SP",
		"ZZZ",
	}

	expected_outs := []*[]Token{
		{
			Token{InstructionToken, I, "ADDI", 0, 4},
			Token{RegisterToken, IGNORE, "X0", 5, 7},
			Token{CommaToken, IGNORE, ",", 7, 8},
			Token{RegisterToken, IGNORE, "X1", 9, 11},
			Token{CommaToken, IGNORE, ",", 11, 12},
			Token{NumberToken, IGNORE, "#12", 13, 16},
		},
		{
			Token{InstructionToken, I, "SUBI", 0, 4},
			Token{RegisterToken, IGNORE, "X1", 5, 7},
			Token{CommaToken, IGNORE, ",", 7, 8},
			Token{RegisterToken, IGNORE, "XZR", 9, 12},
			Token{CommaToken, IGNORE, ",", 12, 13},
			Token{NumberToken, IGNORE, "#9", 14, 16},
		},
		{
			Token{InstructionToken, D, "LDUR", 0, 4},
			Token{RegisterToken, IGNORE, "SP", 5, 7},
			Token{CommaToken, IGNORE, ",", 7, 8},
			Token{LeftBracketToken, IGNORE, "[", 9, 10},
			Token{RegisterToken, IGNORE, "X2", 10, 12},
			Token{CommaToken, IGNORE, ",", 12, 13},
			Token{NumberToken, IGNORE, "#0", 14, 16},
			Token{RightBracketToken, IGNORE, "]", 16, 17},
		},
		{
			Token{InstructionToken, B, "B", 0, 1},
			Token{LabelToken, IGNORE, "label_1", 2, 9},
		},
		{
			Token{InstructionToken, B, "B.EQ", 0, 4},
			Token{LabelToken, IGNORE, "done", 5, 9},
		},
		{
			Token{InstructionToken, B, "BL", 0, 2},
			Token{LabelToken, IGNORE, "link", 3, 7},
		},
		{
			Token{InstructionToken, CB, "CBZ", 0, 3},
			Token{RegisterToken, IGNORE, "X1", 4, 6},
			Token{CommaToken, IGNORE, ",", 6, 7},
			Token{LabelToken, IGNORE, "top", 8, 11},
		},
		{
			Token{InstructionToken, R, "AND", 0, 3},
			Token{RegisterToken, IGNORE, "X12", 4, 7},
			Token{CommaToken, IGNORE, ",", 7, 8},
			Token{RegisterToken, IGNORE, "X10", 9, 12},
			Token{CommaToken, IGNORE, ",", 12, 13},
			Token{RegisterToken, IGNORE, "X1", 14, 16},
		},
		{
			Token{InstructionToken, R, "AND", 0, 3},
			Token{RegisterToken, IGNORE, "X12", 4, 7},
			Token{CommaToken, IGNORE, ",", 7, 8},
			Token{RegisterToken, IGNORE, "X10", 9, 12},
			Token{CommaToken, IGNORE, ",", 12, 13},
			Token{RegisterToken, IGNORE, "XZR", 14, 17},
		},
		{
			Token{InstructionToken, R, "AND", 0, 3},
			Token{RegisterToken, IGNORE, "X12", 4, 7},
			Token{CommaToken, IGNORE, ",", 7, 8},
			Token{RegisterToken, IGNORE, "X10", 9, 12},
			Token{CommaToken, IGNORE, ",", 12, 13},
			Token{RegisterToken, IGNORE, "SP", 14, 16},
		},
		{
			Token{LabelToken, IGNORE, "ZZZ", 0, 3},
		},
	}

	for i, in := range inputs {
		out := TokenizeLine(in)

		if len(*(expected_outs[i])) != len(*out) {
			t.Errorf("Incorrect number of output tokens. Found %d, expected %d. Input=%s. Tokens Received=%v. Tokens Expected=%v", len(*expected_outs[i]), len(*out), in, out, *expected_outs[i])
			continue
		}

		for j, actual := range *out {
			expect := (*(expected_outs)[i])[j]

			if expect.Type != actual.Type {
				t.Errorf("(token=%d) Expected token type %s, got %s. Processing %s. Input %s.", j, expect.Type.String(), actual.Type.String(), actual.Value, in)
			}
			if expect.Value != actual.Value {
				t.Errorf("(token=%d) Expected value %s, got %s. Input %s.", j, expect.Value, actual.Value, in)
			}
			if expect.Start != actual.Start {
				t.Errorf("(token=%d) Expected token start %d, got %d. Input %s.", j, expect.Start, actual.Start, in)
			}
			if expect.End != actual.End {
				t.Errorf("(token=%d) Expected token end %d, got %d. Input %s.", j, expect.End, actual.End, in)
			}

			if j == 0 && expect.Type == InstructionToken && expect.InstructionType != actual.InstructionType {
				t.Errorf("(token=%d) Expected instruction type %s, got %s. Input %s.", j, expect.InstructionType.String(), actual.InstructionType.String(), in)
			}
		}
	}
}
