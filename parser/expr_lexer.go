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
	3, 24715, 42794, 33075, 47597, 16764, 15335, 30598, 22884, 2, 8, 60, 8,
	1, 4, 2, 9, 2, 4, 3, 9, 3, 4, 4, 9, 4, 4, 5, 9, 5, 4, 6, 9, 6, 4, 7, 9,
	7, 4, 8, 9, 8, 4, 9, 9, 9, 4, 10, 9, 10, 4, 11, 9, 11, 4, 12, 9, 12, 3,
	2, 3, 2, 3, 3, 3, 3, 3, 4, 3, 4, 3, 4, 3, 4, 3, 5, 3, 5, 3, 5, 3, 6, 3,
	6, 7, 6, 39, 10, 6, 12, 6, 14, 6, 42, 11, 6, 3, 7, 6, 7, 45, 10, 7, 13,
	7, 14, 7, 46, 3, 7, 3, 7, 3, 8, 3, 8, 3, 9, 3, 9, 3, 10, 3, 10, 3, 11,
	3, 11, 3, 12, 3, 12, 2, 2, 13, 3, 3, 5, 4, 7, 5, 9, 6, 11, 7, 13, 8, 15,
	2, 17, 2, 19, 2, 21, 2, 23, 2, 3, 2, 9, 6, 2, 50, 59, 67, 92, 97, 97, 99,
	124, 5, 2, 11, 12, 15, 15, 34, 34, 4, 2, 67, 67, 99, 99, 4, 2, 80, 80,
	112, 112, 4, 2, 70, 70, 102, 102, 4, 2, 81, 81, 113, 113, 4, 2, 84, 84,
	116, 116, 2, 56, 2, 3, 3, 2, 2, 2, 2, 5, 3, 2, 2, 2, 2, 7, 3, 2, 2, 2,
	2, 9, 3, 2, 2, 2, 2, 11, 3, 2, 2, 2, 2, 13, 3, 2, 2, 2, 3, 25, 3, 2, 2,
	2, 5, 27, 3, 2, 2, 2, 7, 29, 3, 2, 2, 2, 9, 33, 3, 2, 2, 2, 11, 36, 3,
	2, 2, 2, 13, 44, 3, 2, 2, 2, 15, 50, 3, 2, 2, 2, 17, 52, 3, 2, 2, 2, 19,
	54, 3, 2, 2, 2, 21, 56, 3, 2, 2, 2, 23, 58, 3, 2, 2, 2, 25, 26, 7, 42,
	2, 2, 26, 4, 3, 2, 2, 2, 27, 28, 7, 43, 2, 2, 28, 6, 3, 2, 2, 2, 29, 30,
	5, 15, 8, 2, 30, 31, 5, 17, 9, 2, 31, 32, 5, 19, 10, 2, 32, 8, 3, 2, 2,
	2, 33, 34, 5, 21, 11, 2, 34, 35, 5, 23, 12, 2, 35, 10, 3, 2, 2, 2, 36,
	40, 9, 2, 2, 2, 37, 39, 9, 2, 2, 2, 38, 37, 3, 2, 2, 2, 39, 42, 3, 2, 2,
	2, 40, 38, 3, 2, 2, 2, 40, 41, 3, 2, 2, 2, 41, 12, 3, 2, 2, 2, 42, 40,
	3, 2, 2, 2, 43, 45, 9, 3, 2, 2, 44, 43, 3, 2, 2, 2, 45, 46, 3, 2, 2, 2,
	46, 44, 3, 2, 2, 2, 46, 47, 3, 2, 2, 2, 47, 48, 3, 2, 2, 2, 48, 49, 8,
	7, 2, 2, 49, 14, 3, 2, 2, 2, 50, 51, 9, 4, 2, 2, 51, 16, 3, 2, 2, 2, 52,
	53, 9, 5, 2, 2, 53, 18, 3, 2, 2, 2, 54, 55, 9, 6, 2, 2, 55, 20, 3, 2, 2,
	2, 56, 57, 9, 7, 2, 2, 57, 22, 3, 2, 2, 2, 58, 59, 9, 8, 2, 2, 59, 24,
	3, 2, 2, 2, 5, 2, 40, 46, 3, 8, 2, 2,
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
	"", "", "", "AND", "OR", "ID", "WHITESPACE",
}

var lexerRuleNames = []string{
	"T__0", "T__1", "AND", "OR", "ID", "WHITESPACE", "A", "N", "D", "O", "R",
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
	ExprLexerID         = 5
	ExprLexerWHITESPACE = 6
)
