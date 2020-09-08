// Code generated from Expr.g4 by ANTLR 4.8. DO NOT EDIT.

package parser

import (
	"fmt"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = unicode.IsLetter

var serializedLexerAtn = []uint16{
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 9, 75, 8,
	1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9,
	7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 4,
	13, 9, 13, 4, 14, 9, 14, 4, 15, 9, 15, 3, 2, 3, 2, 3, 3, 3, 3, 3, 4, 3,
	4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 6, 3, 6, 3, 6, 3, 6, 3, 6, 3, 7, 3,
	7, 7, 7, 50, 10, 7, 12, 7, 14, 7, 53, 11, 7, 3, 8, 6, 8, 56, 10, 8, 13,
	8, 14, 8, 57, 3, 8, 3, 8, 3, 9, 3, 9, 3, 10, 3, 10, 3, 11, 3, 11, 3, 12,
	3, 12, 3, 13, 3, 13, 3, 14, 3, 14, 3, 15, 3, 15, 2, 2, 16, 3, 3, 5, 4,
	7, 5, 9, 6, 11, 7, 13, 8, 15, 9, 17, 2, 19, 2, 21, 2, 23, 2, 25, 2, 27,
	2, 29, 2, 3, 2, 11, 6, 2, 50, 59, 67, 92, 97, 97, 99, 124, 5, 2, 11, 12,
	15, 15, 34, 34, 4, 2, 67, 67, 99, 99, 4, 2, 80, 80, 112, 112, 4, 2, 70,
	70, 102, 102, 4, 2, 81, 81, 113, 113, 4, 2, 84, 84, 116, 116, 4, 2, 78,
	78, 110, 110, 4, 2, 82, 82, 114, 114, 2, 69, 2, 3, 3, 2, 2, 2, 2, 5, 3,
	2, 2, 2, 2, 7, 3, 2, 2, 2, 2, 9, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13,
	3, 2, 2, 2, 2, 15, 3, 2, 2, 2, 3, 31, 3, 2, 2, 2, 5, 33, 3, 2, 2, 2, 7,
	35, 3, 2, 2, 2, 9, 39, 3, 2, 2, 2, 11, 42, 3, 2, 2, 2, 13, 47, 3, 2, 2,
	2, 15, 55, 3, 2, 2, 2, 17, 61, 3, 2, 2, 2, 19, 63, 3, 2, 2, 2, 21, 65,
	3, 2, 2, 2, 23, 67, 3, 2, 2, 2, 25, 69, 3, 2, 2, 2, 27, 71, 3, 2, 2, 2,
	29, 73, 3, 2, 2, 2, 31, 32, 7, 42, 2, 2, 32, 4, 3, 2, 2, 2, 33, 34, 7,
	43, 2, 2, 34, 6, 3, 2, 2, 2, 35, 36, 5, 17, 9, 2, 36, 37, 5, 19, 10, 2,
	37, 38, 5, 21, 11, 2, 38, 8, 3, 2, 2, 2, 39, 40, 5, 23, 12, 2, 40, 41,
	5, 25, 13, 2, 41, 10, 3, 2, 2, 2, 42, 43, 5, 27, 14, 2, 43, 44, 5, 23,
	12, 2, 44, 45, 5, 23, 12, 2, 45, 46, 5, 29, 15, 2, 46, 12, 3, 2, 2, 2,
	47, 51, 9, 2, 2, 2, 48, 50, 9, 2, 2, 2, 49, 48, 3, 2, 2, 2, 50, 53, 3,
	2, 2, 2, 51, 49, 3, 2, 2, 2, 51, 52, 3, 2, 2, 2, 52, 14, 3, 2, 2, 2, 53,
	51, 3, 2, 2, 2, 54, 56, 9, 3, 2, 2, 55, 54, 3, 2, 2, 2, 56, 57, 3, 2, 2,
	2, 57, 55, 3, 2, 2, 2, 57, 58, 3, 2, 2, 2, 58, 59, 3, 2, 2, 2, 59, 60,
	8, 8, 2, 2, 60, 16, 3, 2, 2, 2, 61, 62, 9, 4, 2, 2, 62, 18, 3, 2, 2, 2,
	63, 64, 9, 5, 2, 2, 64, 20, 3, 2, 2, 2, 65, 66, 9, 6, 2, 2, 66, 22, 3,
	2, 2, 2, 67, 68, 9, 7, 2, 2, 68, 24, 3, 2, 2, 2, 69, 70, 9, 8, 2, 2, 70,
	26, 3, 2, 2, 2, 71, 72, 9, 9, 2, 2, 72, 28, 3, 2, 2, 2, 73, 74, 9, 10,
	2, 2, 74, 30, 3, 2, 2, 2, 5, 2, 51, 57, 3, 8, 2, 2,
}

var lexerDeserializer = antlr.NewATNDeserializer(nil)
var lexerAtn = lexerDeserializer.DeserializeFromUInt16(serializedLexerAtn)

var lexerChannelNames = []string{
	"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
}

var lexerModeNames = []string{
	"DEFAULT_MODE",
}

var lexerLiteralNames = []string{
	"", "'('", "')'",
}

var lexerSymbolicNames = []string{
	"", "", "", "AND", "OR", "LOOP", "ID", "WHITESPACE",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "AND", "OR", "LOOP", "ID", "WHITESPACE", "A", "N", "D",
	"O", "R", "L", "P",
}

type ExprLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var lexerDecisionToDFA = make([]*antlr.DFA, len(lexerAtn.DecisionToState))

func init() {
	for index, ds := range lexerAtn.DecisionToState {
		lexerDecisionToDFA[index] = antlr.NewDFA(ds, index)
	}
}

func NewExprLexer(input antlr.CharStream) *ExprLexer {

	l := new(ExprLexer)

	l.BaseLexer = antlr.NewBaseLexer(input)
	l.Interpreter = antlr.NewLexerATNSimulator(l, lexerAtn, lexerDecisionToDFA, antlr.NewPredictionContextCache())

	l.channelNames = lexerChannelNames
	l.modeNames = lexerModeNames
	l.RuleNames = lexerRuleNames
	l.LiteralNames = lexerLiteralNames
	l.SymbolicNames = lexerSymbolicNames
	l.GrammarFileName = "Expr.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// ExprLexer tokens.
const (
	ExprLexerT__0       = 1
	ExprLexerT__1       = 2
	ExprLexerAND        = 3
	ExprLexerOR         = 4
	ExprLexerLOOP       = 5
	ExprLexerID         = 6
	ExprLexerWHITESPACE = 7
)
