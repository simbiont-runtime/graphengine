// ---

package parser

import (
	"github.com/pingcap/errors"
	"github.com/simbiont-runtime/graphengine/parser/ast"
)

type Parser struct {
	src    string
	lexer  *Lexer
	result []ast.StmtNode

	// the following fields are used by yyParse to reduce allocation.
	cache  []yySymType
	yylval yySymType
	yyVAL  *yySymType
}

func yySetOffset(yyVAL *yySymType, offset int) {
	if yyVAL.expr != nil {
		yyVAL.expr.SetOriginTextPosition(offset)
	}
}

func New() *Parser {
	return &Parser{
		lexer: NewLexer(""),
		cache: make([]yySymType, 200),
	}
}

func (p *Parser) Parse(sql string) (stmts []ast.StmtNode, warns []error, err error) {
	p.lexer.reset(sql)
	p.src = sql
	p.result = p.result[:0]
	yyParse(p.lexer, p)

	warns, errs := p.lexer.Errors()
	if len(warns) > 0 {
		warns = append([]error(nil), warns...)
	} else {
		warns = nil
	}
	if len(errs) != 0 {
		return nil, warns, errors.Trace(errs[0])
	}
	return p.result, warns, nil
}

// ParseOneStmt parses a query and returns an ast.StmtNode.
// The query must have exactly one statement.
func (p *Parser) ParseOneStmt(sql string) (ast.StmtNode, error) {
	stmts, _, err := p.Parse(sql)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if len(stmts) != 1 {
		return nil, errors.New("query must have exactly one statement")
	}
	return stmts[0], nil
}
