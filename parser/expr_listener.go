// Code generated from Expr.g4 by ANTLR 4.8. DO NOT EDIT.

package parser // Expr

import "github.com/antlr/antlr4/runtime/Go/antlr"

// ExprListener is a complete listener for a parse tree produced by ExprParser.
type ExprListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterParenthesis is called when entering the Parenthesis production.
	EnterParenthesis(c *ParenthesisContext)

	// EnterANDOR is called when entering the ANDOR production.
	EnterANDOR(c *ANDORContext)

	// EnterID is called when entering the ID production.
	EnterID(c *IDContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitParenthesis is called when exiting the Parenthesis production.
	ExitParenthesis(c *ParenthesisContext)

	// ExitANDOR is called when exiting the ANDOR production.
	ExitANDOR(c *ANDORContext)

	// ExitID is called when exiting the ID production.
	ExitID(c *IDContext)
}
