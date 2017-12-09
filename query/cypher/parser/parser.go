package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/RossMerr/Caudex.Graph/expressions"
	"github.com/RossMerr/Caudex.Graph/query/cypher/ast"
	"github.com/RossMerr/Caudex.Graph/query/cypher/ir"
	"github.com/RossMerr/Caudex.Graph/query/cypher/lexer"
	"github.com/RossMerr/Caudex.Graph/query/cypher/scanner"
)

// CypherParser represents a parser.
type CypherParser struct {
	s   *scanner.Scanner
	buf struct {
		tok lexer.Type // last read token
		lit string     // last read literal
		n   int        // buffer size (max=1)
	}
}

type Parser interface {
	Parse(r io.Reader) (ast.Clauses, error)
}

var _ Parser = (*CypherParser)(nil)

func (p *CypherParser) label() (string, bool) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != lexer.IDENT && tok == lexer.COLON {
		_, lit = p.scanIgnoreWhitespace()
		return lit, true
	}
	p.unscan()
	return emptyString, false
}

func (p *CypherParser) properties() (map[string]interface{}, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != lexer.IDENT && tok == lexer.LCURLY {

		if properties, err := p.KeyValue(); err == nil {
			tok, lit = p.scanIgnoreWhitespace()
			if tok != lexer.IDENT && tok != lexer.RCURLY {
				return nil, fmt.Errorf("found %q, expected %q", lit, lexer.RCURLY)
			}
			return properties, nil
		}

	}
	p.unscan()
	return nil, nil
}

// KeyValue Loop over all our comma-delimited fields.
func (p *CypherParser) KeyValue() (map[string]interface{}, error) {
	var properties = make(map[string]interface{})
	for {
		tok, lit := p.scanIgnoreWhitespace()
		var prop = lit

		tok, lit = p.scanIgnoreWhitespace()
		if tok != lexer.IDENT && tok != lexer.COLON {
			return nil, fmt.Errorf("found %q, expected %q", lit, lexer.COLON)
		}

		var err error
		if tok, lit, err = p.scanForQuotation(); err == nil {
			properties[prop] = lit
		} else {
			if i, err := strconv.Atoi(lit); err == nil {
				properties[prop] = i
			} else if f, err := strconv.ParseFloat(lit, 64); err == nil {
				properties[prop] = f
			} else if b, err := strconv.ParseBool(lit); err == nil {
				properties[prop] = b
			} else {
				properties[prop] = lit
			}
		}
		tok, _ = p.scanIgnoreWhitespace()
		if tok != lexer.COMMA {
			p.unscan()
			break
		}

	}

	return properties, nil
}

func (p *CypherParser) node() (*ir.VertexPatn, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != lexer.IDENT && tok == lexer.LPAREN {
		stmt := &ir.VertexPatn{}

		tok, lit = p.scanIgnoreWhitespace()
		if tok == lexer.RPAREN {
			return stmt, nil
		} else if tok == lexer.IDENT {
			stmt.Variable = lit
		} else {
			p.unscan()
		}

		if label, ok := p.label(); ok {
			stmt.Label = label
		}

		if properties, err := p.properties(); err == nil && properties != nil {
			stmt.Properties = properties
		} else if err != nil {
			return nil, err
		}

		tok, lit = p.scanIgnoreWhitespace()
		if tok != lexer.IDENT && tok != lexer.RPAREN {
			return nil, fmt.Errorf("found %q, expected %q", lit, lexer.RPAREN)
		}

		return stmt, nil
	}

	p.unscan()
	return nil, nil
}

func (p *CypherParser) length() (uint, uint, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != lexer.IDENT && tok == lexer.MUL {
		min := MinUint
		max := MaxUint

		tok, lit = p.scanIgnoreWhitespace()
		// We have a number
		if tok == lexer.IDENT {
			if u64, err := strconv.ParseUint(lit, 10, 32); err == nil {
				min = uint(u64)
				max = uint(u64)
			} else {
				p.unscan()
			}

			tok, lit = p.scanIgnoreWhitespace()
			if tok == lexer.DOT {
				tok, lit = p.scanIgnoreWhitespace()
				if tok == lexer.DOT {
					max = MaxUint
					tok, lit = p.scanIgnoreWhitespace()
					if u64, err := strconv.ParseUint(lit, 10, 32); err == nil {
						max = uint(u64)
						if min > max {
							return 0, 0, fmt.Errorf("minimum length %d can't exceed maximum length %d for a relationships", min, max)
						}
					} else {
						p.unscan()
					}
				} else {
					return 0, 0, fmt.Errorf("found %q, expected %q", lit, lexer.DOT)
				}
			} else {
				p.unscan()
			}
			// Else we have a range
		} else if tok == lexer.DOT {
			tok, lit = p.scanIgnoreWhitespace()
			if tok == lexer.DOT {
				min = MinUint
				tok, lit = p.scanIgnoreWhitespace()
				if tok == lexer.IDENT {
					if u64, err := strconv.ParseUint(lit, 10, 32); err == nil {
						max = uint(u64)
					} else {
						return 0, 0, fmt.Errorf("found %q, expected uint", lit)
					}
				} else {
					return 0, 0, fmt.Errorf("found %q, expected uint", lit)
				}
			} else {
				return 0, 0, fmt.Errorf("found %q, expected %q", lit, lexer.DOT)
			}
		} else {
			p.unscan()
		}
		return min, max, nil
	}
	p.unscan()

	return 0, 0, nil
}

func (p *CypherParser) relationshipBody() (*ir.EdgeBodyStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != lexer.IDENT && tok == lexer.LSQUARE {
		stmt := &ir.EdgeBodyStmt{}

		tok, lit = p.scanIgnoreWhitespace()
		if tok == lexer.IDENT {
			stmt.Variable = lit
		} else {
			p.unscan()
		}

		if label, ok := p.label(); ok {
			stmt.Type = label
		}

		if min, max, err := p.length(); err == nil && (min != 0 && max != 00) {
			stmt.LengthMinimum = min
			stmt.LengthMaximum = max
		} else if err != nil {
			return nil, err
		} else {
			stmt.LengthMinimum = 1
			stmt.LengthMaximum = 1
		}

		if properties, err := p.properties(); err == nil && properties != nil {
			stmt.Properties = properties
		} else if err != nil {
			return nil, err
		}

		tok, lit := p.scanIgnoreWhitespace()
		if tok != lexer.IDENT && tok != lexer.RSQUARE {
			return nil, fmt.Errorf("found %q, expected %q", lit, lexer.RSQUARE)
		}
		return stmt, nil
	}

	p.unscan()
	return nil, nil
}

func (p *CypherParser) relationship() (*ir.EdgePatn, error) {
	tok, lit := p.scanIgnoreWhitespace()
	// Look for the start of a relationship < or -
	if tok != lexer.IDENT && (tok == lexer.LT || tok == lexer.SUB) {
		stmt := &ir.EdgePatn{Relationship: ir.Undirected}

		if tok == lexer.LT {
			stmt.Relationship = ir.Outbound

			tok, lit = p.scanIgnoreWhitespace()
			// Look for the end of the relationship -
			if tok != lexer.IDENT && tok != lexer.SUB {
				return nil, fmt.Errorf("found %q, expected %q", lit, lexer.SUB)
			}
		}

		if body, err := p.relationshipBody(); err == nil && body != nil {
			stmt.Body = body
		} else if err != nil {
			return nil, err
		}

		tok, lit = p.scanIgnoreWhitespace()
		if tok != lexer.IDENT && tok != lexer.SUB {
			return nil, fmt.Errorf("found %q, expected %q", lit, lexer.SUB)
		}

		// Check for inbound relationship
		if tok == lexer.SUB {
			tok, _ = p.scanIgnoreWhitespace()
			// Look for the end of the relationship - or >
			if tok != lexer.IDENT && tok == lexer.GT {
				stmt.Relationship = ir.Inbound
			} else {
				p.unscan()
			}
		}

		return stmt, nil
	}

	p.unscan()
	return nil, nil
}

func (p *CypherParser) value(tok lexer.Type, lit string) (interface{}, error) {
	//	tok, lit := p.scanIgnoreWhitespace()
	if tok == lexer.SINGLEQUOTATION {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == lexer.IDENT {
			value := lit
			tok, lit := p.scanIgnoreWhitespace()
			if tok == lexer.SINGLEQUOTATION {
				return value, nil
			}

			return emptyString, fmt.Errorf("found %q, expected %q", lit, lexer.SINGLEQUOTATION)
		}

		return emptyString, fmt.Errorf("found %q, expected %q", lit, lexer.IDENT)
	} else if tok == lexer.IDENT {
		if i, err := strconv.Atoi(lit); err == nil {
			return i, nil
		} else if f, err := strconv.ParseFloat(lit, 64); err == nil {
			return f, nil
		} else if b, err := strconv.ParseBool(lit); err == nil {
			return b, nil
		}
		return lit, nil
	}

	p.unscan()
	return emptyString, nil
}

func (p *CypherParser) propertyOrValue() (ast.InterpretExpr, error) {
	tok, lit := p.scanIgnoreWhitespace()

	if tok == lexer.IDENT {

		state := &ast.PropertyStmt{Variable: lit}

		tok2, _ := p.scanIgnoreWhitespace()

		// Must be a value
		if tok2 != lexer.DOT {
			p.unscan()
			value, err := p.value(tok, lit)
			return &ast.Ident{Data: value}, err
		}
		tok, lit = p.scanIgnoreWhitespace()
		if tok != lexer.IDENT {
			return nil, fmt.Errorf("found %q, expected a IDENT", lit)
		}

		state.Value = lit

		return state, nil
	}
	p.unscan()
	return nil, nil
}

func (p *CypherParser) stringValue() ast.InterpretExpr {
	tok, lit, err := p.scanForQuotation()
	if err == nil && tok == lexer.IDENT {
		return &ast.Ident{Data: lit}
	}
	p.unscan()
	return nil
}

func (p *CypherParser) comparisonExpr() (*ast.ComparisonExpr, error) {
	tok, _ := p.scanIgnoreWhitespace()
	switch tok {
	case lexer.EQ:
		return &ast.ComparisonExpr{Comparison: expressions.EQ}, nil
	case lexer.NEQ:
		return &ast.ComparisonExpr{Comparison: expressions.NEQ}, nil
	case lexer.LT:
		return &ast.ComparisonExpr{Comparison: expressions.LT}, nil
	case lexer.LTE:
		return &ast.ComparisonExpr{Comparison: expressions.LTE}, nil
	case lexer.GT:
		return &ast.ComparisonExpr{Comparison: expressions.GT}, nil
	case lexer.GTE:
		return &ast.ComparisonExpr{Comparison: expressions.GTE}, nil
	}
	p.unscan()
	return nil, nil
}

func (p *CypherParser) booleanExpr() (ast.InterpretExpr, error) {
	tok, _ := p.scanIgnoreWhitespace()
	switch tok {
	case lexer.AND:
		return &ast.BooleanExpr{Boolean: expressions.AND}, nil
	case lexer.OR:
		return &ast.BooleanExpr{Boolean: expressions.OR}, nil
	case lexer.NOT:
		return &ast.NotExpr{}, nil
	case lexer.XOR:
		return &ast.BooleanExpr{Boolean: expressions.XOR}, nil
	}
	p.unscan()
	return nil, nil
}

// Predicate pulls of each item to pass into the shunting algorithm to build up the AST
func (p *CypherParser) Predicate() (ast.Expr, error) {
	exprStack := make(StackExpr, 0)

	tok, _ := p.scanIgnoreWhitespace()
	p.unscan()

	for !tok.IsClause() && tok != lexer.EOF {

		if property, err := p.propertyOrValue(); err == nil && property != nil {
			exprStack = exprStack.Push(property)
		} else if err != nil {
			return nil, err
		} else if property := p.stringValue(); property != nil {
			exprStack = exprStack.Push(property)
		} else if err != nil {
			return nil, err
		} else if comparisonExpr, err := p.comparisonExpr(); err == nil && comparisonExpr != nil {
			exprStack = exprStack.Push(comparisonExpr)
		} else if err != nil {
			return nil, err
		} else if booleanExpr, err := p.booleanExpr(); err == nil && booleanExpr != nil {
			exprStack = exprStack.Push(booleanExpr)
		} else if err != nil {
			return nil, err
		}

		tok, _ = p.scanIgnoreWhitespace()
		p.unscan()
	}

	root, err := exprStack.Shunt()

	return root, err
}

func (p *CypherParser) where() (*ast.WhereStmt, error) {
	tok, _ := p.scanIgnoreWhitespace()
	if tok == lexer.WHERE {
		state := &ast.WhereStmt{}

		if predicate, err := p.Predicate(); err == nil {
			state.Predicate = predicate
		} else {
			return nil, err
		}

		return state, nil
	}

	p.unscan()
	return nil, nil
}

func (p *CypherParser) match() (ast.Clauses, ast.Clauses, error) {
	state := &ast.MatchStmt{}
	pattern, err := p.pattern()
	if err == nil {
		state.Pattern = pattern

		where, err := p.where()
		if err == nil && where != nil {
			state.Next = where
		} else if err != nil {
			return nil, nil, err
		}

		returns, err := p.returns()
		if err != nil {
			return nil, nil, err
		}

		if where != nil {
			where.Next = returns
		} else {
			state.Next = returns

		}

		return state, returns, nil
	}
	return nil, nil, err
}

func (p *CypherParser) optionalMatch() (ast.Clauses, error) {
	state := &ast.OptionalMatchStmt{}
	pattern, err := p.pattern()
	if err == nil {
		state.Pattern = pattern
		//state.Next = next
		return state, nil
	}
	return nil, err
}

func (p *CypherParser) pattern() (ir.Patn, error) {
	var pattern ir.Patn
	var lastVertex *ir.VertexPatn
	var lastEdge *ir.EdgePatn

	// Next we should loop over all the pattern.
	for {

		if node, err := p.node(); err == nil && node != nil {
			lastVertex = node
			if pattern == nil {
				pattern = lastVertex
			}
			if lastEdge != nil {
				lastEdge.Vertex = node
			}
		} else if err != nil {
			return nil, err
		}

		if relationship, err := p.relationship(); err == nil && relationship != nil {
			lastEdge = relationship
			lastVertex.Edge = relationship
		} else if err != nil {
			return nil, err
		} else {
			break
		}
	}

	return pattern, nil
}

func (p *CypherParser) returns() (*ast.ReturnStmt, error) {
	tok, _ := p.scanIgnoreWhitespace()
	if tok == lexer.RETURN {
		state := &ast.ReturnStmt{}

		if maps, err := p.MapVariables(); err == nil {
			state.Maps = maps
		} else {
			return nil, err
		}

		return state, nil
	}

	p.unscan()
	return nil, nil
}

func (p *CypherParser) MapAlias() (string, error) {
	tok, _ := p.scanIgnoreWhitespace()
	if tok == lexer.AS {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == lexer.IDENT {
			return lit, nil
		}
		return "", fmt.Errorf("found %q, expected column alias", lit)
	}

	p.unscan()
	return "", nil
}

func (p *CypherParser) MapElement() (ast.MapElementStmt, error) {
	tok, lit := p.scanIgnoreWhitespace()

	if tok == lexer.DOT {
		tok, lit := p.scanIgnoreWhitespace()

		if tok == lexer.IDENT {
			alias, err := p.MapAlias()
			if err != nil {
				return nil, err
			}

			return &ast.MapProperty{Key: lit, Alias: alias}, nil
		} else if tok == lexer.MUL {
			return &ast.MapAll{}, nil
		} else {
			return nil, fmt.Errorf("found %q, expected part of a map", lit)
		}
	} else if tok == lexer.IDENT {
		key := lit

		tok, _ := p.scanIgnoreWhitespace()

		if tok == lexer.COLON {
			return nil, fmt.Errorf("found %q, MapLiteral not yet supported", lit)
			// todo MapLiteral
			// literal := &ast.MapLiteral{Key: key}
			// mapPro.Elements = append(mapPro.Elements, literal)
		} else {
			p.unscan()
			return &ast.MapVariable{Key: key}, nil
		}
	}
	p.unscan()
	return nil, nil
}

func (p *CypherParser) MapElements() ([]ast.MapElementStmt, error) {

	elements := make([]ast.MapElementStmt, 0)

	for {
		if e, err := p.MapElement(); e != nil && err == nil {
			elements = append(elements, e)
		}

		tok, lit := p.scanIgnoreWhitespace()

		if tok == lexer.COMMA {
			continue
		} else if tok == lexer.RCURLY {
			p.unscan()
			break
		} else {
			return nil, fmt.Errorf("found %q, expected } or ,", lit)
		}
	}

	return elements, nil
}
func (p *CypherParser) MapVariables() ([]*ast.MapProjectionStmt, error) {
	maps := make(map[string]*ast.MapProjectionStmt)
	for {
		tok, lit, err := p.scanForQuotation()

		if tok == lexer.GRAVE {
			p.unscan()
			if tok, lit, err = p.scanForQuotation(); err != nil {
				return nil, err
			}
			if tok.IsQuotation() {
				tok = lexer.IDENT
			}

		}

		if tok == lexer.MUL {
			if _, ok := maps[lit]; !ok {
				maps["*"] = ast.NewMapProjectionStmt("*", &ast.MapAll{})
			}
		}

		if tok == lexer.IDENT {

			if _, ok := maps[lit]; !ok {
				maps[lit] = ast.NewMapProjectionStmt(lit)
			}

			tok, _ := p.scanIgnoreWhitespace()

			if tok == lexer.LCURLY {

				if elements, err := p.MapElements(); err == nil && elements != nil {
					maps[lit].Elements = append(maps[lit].Elements, elements...)
				}

				tok, _ := p.scanIgnoreWhitespace()

				if tok != lexer.RCURLY {
					return nil, fmt.Errorf("found %q, expected }", lit)
				}
			} else if tok == lexer.DOT {
				p.unscan()
				if element, err := p.MapElement(); err == nil && element != nil {
					maps[lit].Elements = append(maps[lit].Elements, element)
				}
			}

		}

		tok, _ = p.scanIgnoreWhitespace()

		if tok != lexer.COMMA {
			p.unscan()
			break
		}

	}

	arr := make([]*ast.MapProjectionStmt, 0)
	for _, m := range maps {
		if len(m.Elements) == 0 {
			m.Elements = append(m.Elements, &ast.MapAll{})
		}
		arr = append(arr, m)
	}

	return arr, nil
}

func (p *CypherParser) create() (ast.Clauses, error) {
	state := &ast.CreateStmt{}
	pattern, err := p.pattern()
	if err == nil {
		state.Pattern = pattern
		//	state.Next = next
		return state, nil
	}
	return nil, err
}

func (p *CypherParser) delete() (ast.Clauses, error) {
	state := &ast.DeleteStmt{}
	pattern, err := p.pattern()
	if err == nil {
		state.Pattern = pattern
		//	state.Next = next
		return state, nil
	}
	return nil, err
}

func (p *CypherParser) clause() (ast.Clauses, error) {
	tok, lit := p.scanIgnoreWhitespace()

	if !tok.IsClause() {
		return nil, fmt.Errorf("found %q, expected a clause", lit)
	}

	if tok == lexer.OPTIONAL {
		tok, lit = p.scanIgnoreWhitespace()
		if tok == lexer.MATCH {
			tok = lexer.OPTIONAL_MATCH
		} else {
			return nil, fmt.Errorf("found %q, expected MATCH", lit)
		}
	} else if tok == lexer.DETACH {
		tok, lit = p.scanIgnoreWhitespace()
		if tok == lexer.DELETE {
			tok = lexer.DETACH_DELETE
		} else {
			return nil, fmt.Errorf("found %q, expected DELETE", lit)
		}
	}

	switch tok {
	case lexer.MATCH:
		root, _, err := p.match()
		return root, err
	case lexer.OPTIONAL_MATCH:
		return p.optionalMatch()
	case lexer.CREATE:
		return p.create()
	case lexer.DELETE:
		return p.delete()
	}

	return nil, fmt.Errorf("No matching statement found %q", lit)
}

func (p *CypherParser) subClause() (lexer.Type, bool) {
	tok, _ := p.scanIgnoreWhitespace()

	if tok.IsSubClause() {
		if tok == lexer.ON {
			tok, lit := p.scanIgnoreWhitespace()
			if tok == lexer.CREATE {
				return lexer.ON_CREATE, true
			} else if tok == lexer.MATCH {
				return lexer.ON_MATCH, true
			} else {
				panic(fmt.Sprintf("found %q, expected CREATE", lit))
			}
		}

		if tok == lexer.ORDER {
			tok, lit := p.scanIgnoreWhitespace()
			if tok == lexer.BY {
				return lexer.ORDER_BY, true
			}
			panic(fmt.Sprintf("found %q, expected BY", lit))
		}

		return tok, true
	}

	p.unscan()
	return lexer.IDENT, false
}

// Parse parses a cypher Clauses statement.
func (p *CypherParser) Parse(r io.Reader) (ast.Clauses, error) {
	p.s = scanner.NewScanner(r)
	return p.clause()
}

// NewParser returns a new instance of Parser.
func NewParser() *CypherParser {
	return &CypherParser{}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *CypherParser) scan() (tok lexer.Type, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace lexer.
func (p *CypherParser) scanIgnoreWhitespace() (tok lexer.Type, lit string) {
	tok, lit = p.scan()
	if tok == lexer.WS {
		tok, lit = p.scan()
	}
	return
}

// scanForQuotation scans the next matching quotations lexer.
func (p *CypherParser) scanForQuotation() (tok lexer.Type, lit string, err error) {
	tok, lit = p.scanIgnoreWhitespace()
	if tok == lexer.QUOTATION || tok == lexer.SINGLEQUOTATION || tok == lexer.GRAVE {
		lit = emptyString
		for {
			tok2, s := p.scan()
			if tok2 != tok {
				lit += s
			} else if tok2 == lexer.EOF {
				err = fmt.Errorf("No matching quotaation found %q", lit)
				return
			} else {
				tok = lexer.IDENT
				break
			}
		}
		return
	}
	err = fmt.Errorf("No matching quotaation found %q", lit)
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *CypherParser) unscan() { p.buf.n = 1 }
