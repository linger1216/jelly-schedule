package core

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/linger1216/jelly-schedule/parser"
)

// func (e *Executor) getJob(jobId string) (Job, error) {
type ExprListener struct {
	err   error
	stack []Job
	getFn func(string) (Job, error)
	andFn func(left, right Job) Job
	orFn  func(left, right Job) Job
}

func NewExprListener(
	getFn func(string) (Job, error),
	andFn func(left, right Job) Job,
	orFn func(left, right Job) Job) *ExprListener {
	return &ExprListener{err: nil, getFn: getFn, andFn: andFn, orFn: orFn}
}

func (e *ExprListener) error() error {
	return e.err
}

func (e *ExprListener) push(i Job) {
	e.stack = append(e.stack, i)
}

func (e *ExprListener) Pop() Job {
	if len(e.stack) < 1 {
		panic("stack is empty unable to Pop")
	}
	result := e.stack[len(e.stack)-1]
	e.stack = e.stack[:len(e.stack)-1]
	return result
}

func (e *ExprListener) VisitTerminal(node antlr.TerminalNode) {
}

func (e *ExprListener) VisitErrorNode(node antlr.ErrorNode) {
}

func (e *ExprListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
}

func (e *ExprListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
}

func (e *ExprListener) EnterStart(c *parser.StartContext) {
}

func (e *ExprListener) EnterParenthesis(c *parser.ParenthesisContext) {
}

func (e *ExprListener) EnterANDOR(c *parser.ANDORContext) {
}

func (e *ExprListener) EnterID(c *parser.IDContext) {
}

func (e *ExprListener) ExitStart(c *parser.StartContext) {
}

func (e *ExprListener) ExitParenthesis(c *parser.ParenthesisContext) {
}

func (e *ExprListener) ExitANDOR(c *parser.ANDORContext) {
	if e.err != nil {
		return
	}

	right, left := e.Pop(), e.Pop()
	//_MOD(_Expr).Debugf("parser left:%s right:%s", left, right)
	switch c.GetOp().GetTokenType() {
	case parser.ExprLexerAND:
		e.push(e.andFn(left, right))
	case parser.ExprLexerOR:
		e.push(e.orFn(left, right))
	default:
		panic(fmt.Sprintf("unexpected op: %s", c.GetOp().GetText()))
	}
}

func (e *ExprListener) ExitID(c *parser.IDContext) {
	if e.err != nil {
		return
	}

	job, err := e.getFn(c.GetText())
	if err != nil {
		_MOD(_Expr).With(_Job, c.GetText()).Debugf("not found")
		e.err = err
		return
	}

	if job != nil {
		e.push(job)
	}
}
