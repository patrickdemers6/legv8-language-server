package languageserver

var KeywordInstructionTypes map[string]InstructionType

type TokenType int8

type InstructionType int8

const (
	NONE InstructionType = iota
	I
	R
	D
	B
	BR
	CB
	IM
	IGNORE
	UNKNOWN
)

func (instruction InstructionType) Check(line string, current int) *Token {
	l := isInstructionType(line, current, instruction)

	if l > 0 {
		return &Token{
			Value:           line[current : current+l],
			Type:            InstructionToken,
			InstructionType: instruction,
			Start:           current,
			End:             current + l,
		}
	}
	return nil
}

const (
	RegisterToken TokenType = iota
	CommaToken
	InstructionToken
	LabelToken
	LeftBracketToken
	RightBracketToken
	UnknownToken
	EOLToken
	NumberToken
	ColonToken
)

func (inst InstructionType) String() string {
	switch inst {
	case I:
		return "I"
	case R:
		return "R"
	case D:
		return "D"
	case B:
		return "B"
	case BR:
		return "BR"
	case CB:
		return "CB"
	case IM:
		return "IM"
	case IGNORE:
		return "IGNORE"
	}
	return "UNKNOWN"
}

func (t TokenType) String() string {
	switch t {
	case CommaToken:
		return "Comma"
	case RegisterToken:
		return "Register"
	case InstructionToken:
		return "Instruction"
	case LabelToken:
		return "Label"
	case EOLToken:
		return "EOL"
	case LeftBracketToken:
		return "Left Bracket"
	case RightBracketToken:
		return "Right Bracket"
	case NumberToken:
		return "Immediate"
	case ColonToken:
		return "Colon"
	}
	return "Unknown"
}

type Token struct {
	Type            TokenType
	InstructionType InstructionType
	Value           string
	Start           int
	End             int
}

func init() {
	KeywordInstructionTypes = make(map[string]InstructionType)

	iTypeInstructions := []string{"ADDI", "SUBI", "ANDI", "ADDIS", "ORRI", "EORI", "SUBIS", "ANDIS", "LSL", "LSR"}
	imTypeInstructions := []string{"PRNT"}
	dTypeInstructions := []string{"STURB", "LDURB", "STURH", "LDURH", "STURW", "LDURSW", "STXR", "LDXR", "STUR", "LDUR"}
	rTypeInstructions := []string{"FDIVS", "FMULS", "FCMPS", "FADDS", "FSUBS", "FMULD", "FDIVD", "FCMPD", "FADDD", "FSUBD", "AND", "ADD", "SDIV", "UDIV", "MUL", "SMULH", "UMULH", "ORR", "ADDS", "STURS", "LDURS", "EOR", "SUB", "ANDS", "SUBS", "STURD", "LDURD"}
	bTypeInstructions := []string{"B.EQ", "B.GT", "B.NE", "B.HS", "B.LO", "B.MI", "B.PL", "B.VS", "B.VC", "B.HI", "B.LS", "B.GE", "B.LT", "B.LE", "B", "BL"}
	cbTypeInstructions := []string{"CBZ", "CBNZ"}
	brTypeInstructions := []string{"BR"}
	ignoreTypeInstructions := []string{"PRNL", "DUMP", "HALT"}

	for _, v := range iTypeInstructions {
		KeywordInstructionTypes[v] = I
	}
	for _, v := range imTypeInstructions {
		KeywordInstructionTypes[v] = IM
	}
	for _, v := range dTypeInstructions {
		KeywordInstructionTypes[v] = D
	}
	for _, v := range rTypeInstructions {
		KeywordInstructionTypes[v] = R
	}
	for _, v := range bTypeInstructions {
		KeywordInstructionTypes[v] = B
	}
	for _, v := range cbTypeInstructions {
		KeywordInstructionTypes[v] = CB
	}
	for _, v := range brTypeInstructions {
		KeywordInstructionTypes[v] = BR
	}
	for _, v := range ignoreTypeInstructions {
		KeywordInstructionTypes[v] = IGNORE
	}
}
