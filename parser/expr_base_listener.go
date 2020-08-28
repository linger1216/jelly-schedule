// Code generated from Expr.g4 by ANTLR 4.8. DO NOT EDIT.

package parser // Expr

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseExprListener is a complete listener for a parse tree produced by ExprParser.
type BaseExprListener struct{}

var _ ExprListener = &BaseExprListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseExprListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseExprListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseExprListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseExprListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseExprListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseExprListener) ExitStart(ctx *StartContext) {}

// EnterParenthesis is called when production Parenthesis is entered.
func (s *BaseExprListener) EnterParenthesis(ctx *ParenthesisContext) {}

// ExitParenthesis is called when production Parenthesis is exited.
func (s *BaseExprListener) ExitParenthesis(ctx *ParenthesisContext) {}

// EnterANDOR is called when production ANDOR is entered.
func (s *BaseExprListener) EnterANDOR(ctx *ANDORContext) {}

// ExitANDOR is called when production ANDOR is exited.
func (s *BaseExprListener) ExitANDOR(ctx *ANDORContext) {}

// EnterID is called when production ID is entered.
func (s *BaseExprListener) EnterID(ctx *IDContext) {}

// ExitID is called when production ID is exited.
func (s *BaseExprListener) ExitID(ctx *IDContext) {}
