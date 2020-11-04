package token

import (
	"github.com/chronos-tachyon/go-spiderscript/internal/util"
	"github.com/chronos-tachyon/go-spiderscript/value"
)

type Lexer struct {
	inp []rune
	tok Token
	pos Position
	st  State
	ctr int32
	rdy bool
	eof bool
}

func NewLexer(path string, input []rune) *Lexer {
	lexer := new(Lexer)
	lexer.Init(path, input)
	return lexer
}

func (lexer *Lexer) Init(path string, input []rune) {
	*lexer = Lexer{
		inp: input,
		pos: Position{Path: path},
		st:  StateReadyAtStartOfFile,
	}
}

func (lexer *Lexer) HasNext() bool {
	if lexer.rdy {
		return true
	}

	if lexer.eof {
		return false
	}

	if lexer.st != StateReady && lexer.st != StateReadyAtStartOfFile {
		lexer.st = StateReady
	}

	lexer.tok = Token{
		Type:  Partial,
		Start: lexer.pos,
		End:   lexer.pos,
	}

	for lexer.pos.RuneOffset < uint(len(lexer.inp)) {
		ch := lexer.inp[lexer.pos.RuneOffset]
		consumeRune, isComplete := lexer.processRune(ch)
		if consumeRune {
			lexer.pos.Advance(ch)
		}
		if isComplete {
			break
		}
	}

	lexer.tok.End = lexer.pos

	if lexer.tok.Start.RuneOffset >= lexer.pos.RuneOffset {
		lexer.tok.Type = EOF
		util.EstimateLengths(&lexer.tok)
		lexer.rdy = true
		lexer.eof = true
		return true
	}

	lexer.tok.Raw = lexer.inp[lexer.tok.Start.RuneOffset:lexer.pos.RuneOffset]
	rawStr := string(lexer.tok.Raw)

	if lexer.tok.Type == Partial {
		if tt, ok := earlyEOFTable[lexer.st]; ok {
			lexer.tok.Type = tt
		}
	}

	if lexer.tok.Type == Identifier {
		if tt, ok := keywordTable[rawStr]; ok {
			lexer.tok.Type = tt
		}
	}

	switch lexer.tok.Type {
	case Invalid, Partial, ShebangLine, SingleLineComment, MultiLineComment, Pragma, Identifier:
		lexer.tok.Parsed = &value.Literal{rawStr}

	case Number:
		nv := &value.Number{}
		if err := nv.Parse(lexer.tok.Raw); err != nil {
			ev := &value.Error{err, 32, 64}
			lexer.tok.Parsed = ev
		} else {
			lexer.tok.Parsed = nv
		}

	case String:
		sv := &value.String{}
		if err := sv.Parse(lexer.tok.Raw); err != nil {
			ev := &value.Error{err, 32, 64}
			lexer.tok.Parsed = ev
		} else {
			lexer.tok.Parsed = sv
		}

	case Regex:
		rx := &value.Regex{}
		if err := rx.Parse(lexer.tok.Raw); err != nil {
			ev := &value.Error{err, 32, 64}
			lexer.tok.Parsed = ev
		} else {
			lexer.tok.Parsed = rx
		}

	case PEG:
		pv := &value.PEG{}
		if err := pv.Parse(lexer.tok.Raw); err != nil {
			ev := &value.Error{err, 32, 64}
			lexer.tok.Parsed = ev
		} else {
			lexer.tok.Parsed = pv
		}

	case HWS:
		p := lexer.tok.Start.RawColumnNumber
		q := lexer.tok.End.RawColumnNumber
		if p > q {
			p = q
		}
		tmp := make([]byte, q-p)
		for i := range tmp {
			tmp[i] = ' '
		}
		lexer.tok.Parsed = value.NewLiteral(string(tmp))

	case VWS:
		p := lexer.tok.Start.RawLineNumber
		q := lexer.tok.End.RawLineNumber
		if p > q {
			p = q
		}
		tmp := make([]byte, q-p)
		for i := range tmp {
			tmp[i] = '\n'
		}
		lexer.tok.Parsed = value.NewLiteral(string(tmp))
	}

	if est, ok := lexer.tok.Parsed.(util.Estimable); ok {
		util.EstimateLengths(est)
	}
	util.EstimateLengths(&lexer.tok)
	lexer.rdy = true
	return true
}

func (lexer *Lexer) Token() Token {
	if lexer.rdy {
		lexer.rdy = false
		return lexer.tok
	}

	tok := Token{
		Type:  Invalid,
		Start: lexer.pos,
		End:   lexer.pos,
	}
	return tok
}

func (lexer *Lexer) processRune(ch rune) (consumeRune bool, isComplete bool) {
	for _, row := range stateTable {
		newCounter := lexer.ctr + int32(row.Delta)

		if !row.CurrentState.MatchState(lexer.st) {
			continue
		}

		if !row.CurrentRune.MatchRune(ch) {
			continue
		}

		if row.CurrentRaw != nil {
			raw := lexer.inp[lexer.tok.Start.RuneOffset:lexer.pos.RuneOffset]
			if !util.EqualRunes(raw, row.CurrentRaw) {
				continue
			}
		}

		if row.OnlyIfZero && newCounter != 0 {
			continue
		}

		if row.NextState == StateZero {
			panic("invalid use of State(0)")
		}

		if row.NextState == StateDone {
			lexer.tok.Type = row.TokenType
			lexer.ctr = 0
			return true, true
		}

		if row.NextState == StateUnreadAndDone {
			lexer.tok.Type = row.TokenType
			lexer.ctr = 0
			return false, true
		}

		if row.NextState == StateIllegalRune {
			break
		}

		lexer.st = row.NextState
		lexer.ctr = newCounter
		return true, false
	}

	if lexer.pos.RuneOffset > lexer.tok.Start.RuneOffset {
		lexer.tok.Type = Partial
		lexer.ctr = 0
		return false, true
	}

	lexer.tok.Type = Invalid
	lexer.ctr = 0
	return true, true
}
