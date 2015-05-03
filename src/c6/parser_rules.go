package c6

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

import "fmt"
import "strconv"
import "c6/ast"

func (parser *Parser) ParseStatement(parentRuleSet *ast.RuleSet) ast.Statement {
	var token = parser.peek()

	if token.Type == ast.T_IMPORT {
		return parser.ParseImportStatement()
	} else if token.Type == ast.T_VARIABLE {
		return parser.ParseVariableAssignment()
	} else if token.IsSelector() {
		return parser.ParseRuleSet(parentRuleSet)
	}
	return nil
}

func (parser *Parser) ParseRuleSet(parentRuleSet *ast.RuleSet) ast.Statement {
	var ruleset = ast.RuleSet{}
	var tok = parser.next()

	for tok.IsSelector() {

		switch tok.Type {

		case ast.T_TYPE_SELECTOR:
			sel := ast.TypeSelector{tok.Str}
			ruleset.AppendSelector(sel)

		case ast.T_UNIVERSAL_SELECTOR:
			sel := ast.UniversalSelector{}
			ruleset.AppendSelector(sel)

		case ast.T_ID_SELECTOR:
			sel := ast.IdSelector{tok.Str}
			ruleset.AppendSelector(sel)

		case ast.T_CLASS_SELECTOR:
			sel := ast.ClassSelector{tok.Str}
			ruleset.AppendSelector(sel)

		case ast.T_PARENT_SELECTOR:
			sel := ast.ParentSelector{parentRuleSet}
			ruleset.AppendSelector(sel)

		case ast.T_PSEUDO_SELECTOR:
			sel := ast.PseudoSelector{tok.Str, ""}
			if nextTok := parser.peek(); nextTok.Type == ast.T_LANG_CODE {
				sel.C = nextTok.Str
			}
			ruleset.AppendSelector(sel)

		case ast.T_ADJACENT_SELECTOR:
			ruleset.AppendSelector(ast.AdjacentSelector{})
		case ast.T_CHILD_SELECTOR:
			ruleset.AppendSelector(ast.ChildSelector{})
		case ast.T_DESCENDANT_SELECTOR:
			ruleset.AppendSelector(ast.DescendantSelector{})
		default:
			panic(fmt.Errorf("Unexpected selector token: %+v", tok))
		}
		tok = parser.next()
	}
	parser.backup()

	// parse declaration block
	ruleset.DeclarationBlock = parser.ParseDeclarationBlock(&ruleset)
	return &ruleset
}

/**
This method returns objects with ast.Number interface

works for:

	'10'
	'10' 'px'
	'10' 'em'
	'0.2' 'em'
*/
func (parser *Parser) ParseNumber() ast.Expression {
	var pos = parser.Pos
	debug("ParseNumber at %d", parser.Pos)

	// the number token
	var tok = parser.next()
	debug("ParseNumber => next: %s", tok)

	var negative = false

	if tok.Type == ast.T_MINUS {
		tok = parser.next()
		negative = true
	} else if tok.Type == ast.T_PLUS {
		tok = parser.next()
		negative = false
	}

	var val float64
	var tok2 = parser.peek()

	if tok.Type == ast.T_INTEGER {

		i, err := strconv.ParseInt(tok.Str, 10, 64)
		if err != nil {
			panic(err)
		}
		if negative {
			i = -i
		}
		val = float64(i)

	} else if tok.Type == ast.T_FLOAT {

		f, err := strconv.ParseFloat(tok.Str, 64)
		if err != nil {
			panic(err)
		}
		if negative {
			f = -f
		}
		val = f

	} else {
		// unknown token
		parser.restore(pos)
		return nil
	}

	if tok2.IsOneOfTypes([]ast.TokenType{ast.T_UNIT_PX, ast.T_UNIT_PT, ast.T_UNIT_CM, ast.T_UNIT_EM, ast.T_UNIT_MM, ast.T_UNIT_REM, ast.T_UNIT_DEG, ast.T_UNIT_PERCENT}) {
		// consume the unit token
		parser.next()
		return ast.NewLength(val, ast.ConvertTokenTypeToUnitType(tok2.Type), tok)
	}
	return ast.NewNumber(val, tok)
}

func (parser *Parser) ParseFunctionCall() *ast.FunctionCall {
	var identTok = parser.next()

	debug("ParseFunctionCall => next: %s", identTok)

	var fcall = ast.NewFunctionCall(identTok)

	parser.expect(ast.T_PAREN_START)

	var argTok = parser.peek()
	for argTok.Type != ast.T_PAREN_END {
		var arg = parser.ParseFactor()
		fcall.AppendArgument(arg)
		debug("ParseFunctionCall => arg: %+v", arg)

		argTok = parser.peek()
		if argTok.Type == ast.T_COMMA {
			parser.next() // skip comma
			argTok = parser.peek()
		} else if argTok.Type == ast.T_PAREN_END {
			parser.next() // consume ')'
			break
		}
	}
	return fcall
}

func (parser *Parser) ParseIdent() *ast.Ident {
	var tok = parser.next()
	debug("ReduceIndent => next: %s", tok)
	if tok.Type != ast.T_IDENT {
		panic("Invalid token for ident.")
	}
	return ast.NewIdent(tok.Str, *tok)
}

/**
The ParseFactor must return an Expression interface compatible object
*/
func (parser *Parser) ParseFactor() ast.Expression {
	debug("ParseFactor at %d", parser.Pos)
	var tok = parser.peek()
	debug("ParseFactor => peek: %s", tok)

	if tok.Type == ast.T_PAREN_START {

		parser.expect(ast.T_PAREN_START)
		var expr = parser.ParseExpression()
		parser.expect(ast.T_PAREN_END)
		return expr

	} else if tok.Type == ast.T_INTERPOLATION_START {

		return parser.ParseInterp()

	} else if tok.Type == ast.T_QQ_STRING {

		tok = parser.next()
		var str = ast.NewStringWithQuote('"', tok)
		return ast.Expression(str)

	} else if tok.Type == ast.T_Q_STRING {

		tok = parser.next()
		var str = ast.NewStringWithQuote('\'', tok)
		return ast.Expression(str)

	} else if tok.Type == ast.T_IDENT {

		tok = parser.next()
		return ast.Expression(ast.NewString(tok))

	} else if tok.Type == ast.T_INTEGER || tok.Type == ast.T_FLOAT {

		// reduce number
		var number = parser.ParseNumber()
		return ast.Expression(number)

	} else if tok.Type == ast.T_FUNCTION_NAME {

		var fcall = parser.ParseFunctionCall()
		return ast.Expression(*fcall)

	} else if tok.Type == ast.T_HEX_COLOR {

		panic("hex color is not implemented yet")

		// TODO: Add more incorrect cases here

	} else {

		return nil
	}
	return nil
}

func (parser *Parser) ParseTerm() ast.Expression {
	debug("ParseTerm at %d", parser.Pos)
	var pos = parser.Pos
	var factor = parser.ParseFactor()
	if factor == nil {
		parser.restore(pos)
		return nil
	}

	// see if the next token is '*' or '/'
	var tok = parser.peek()
	if tok.Type == ast.T_MUL || tok.Type == ast.T_DIV {
		parser.next()
		if term := parser.ParseTerm(); term != nil {
			if tok.Type == ast.T_MUL {
				return ast.NewBinaryExpression(ast.OpMul, factor, term)
			} else if tok.Type == ast.T_DIV {
				return ast.NewBinaryExpression(ast.OpDiv, factor, term)
			}
		} else {
			panic("Unexpected token after * and /")
		}
	}
	return factor
}

/**

We here treat the property values as expressions:

	padding: {expression} {expression} {expression};
	margin: {expression};

*/
func (parser *Parser) ParseExpression() ast.Expression {
	var pos = parser.Pos

	debug("ParseExpression")

	// plus or minus. this creates an unary expression that holds the later term.
	// this is for:  +3 or -4
	var tok = parser.peek()
	var expr ast.Expression
	if tok.Type == ast.T_PLUS || tok.Type == ast.T_MINUS {
		parser.next()
		if term := parser.ParseTerm(); term != nil {
			expr = ast.NewUnaryExpression(ast.ConvertTokenTypeToOpType(tok.Type), term)
		} else {
			parser.restore(pos)
			return nil
		}
	} else {
		expr = parser.ParseTerm()
	}

	if expr == nil {
		debug("ParseExpression failed, got %+v, restoring to %d", expr, pos)
		parser.restore(pos)
		return nil
	}

	var rightTok = parser.peek()
	for rightTok.Type == ast.T_PLUS || rightTok.Type == ast.T_MINUS || rightTok.Type == ast.T_LITERAL_CONCAT {
		// accept plus or minus
		parser.next()

		if rightTerm := parser.ParseTerm(); rightTerm != nil {
			expr = ast.NewBinaryExpression(ast.ConvertTokenTypeToOpType(rightTok.Type), expr, rightTerm)
		} else {
			panic("right term is not parseable.")
		}
		rightTok = parser.peek()
	}
	return expr
}

/*
func (parser *Parser) ParseMap() *ast.Map {
	var tok = parser.next()
	if tok.Type != ast.T_PAREN_START {
		panic("Map Syntax error: expecting '('")
	}

	tok = parser.next()
	if tok.Type != ast.T_IDENT {
		panic("Map Syntax error: expecting ident or expression after '('")
	}
	return nil
}
*/

func (parser *Parser) ParseMap() ast.Expression {
	var pos = parser.Pos
	var tok = parser.next()
	// since it's not started with '(', it's not map
	if tok.Type != ast.T_PAREN_START {
		parser.restore(pos)
		return nil
	}

	tok = parser.peek()
	for tok.Type != ast.T_PAREN_END {
		var keyExpr = parser.ParseExpression()
		if keyExpr == nil {
			parser.restore(pos)
			return nil
		}

		if parser.expect(ast.T_COLON) == nil {
			parser.restore(pos)
			return nil
		}

		var valueExpr = parser.ParseExpression()
		if valueExpr == nil {
			parser.restore(pos)
			return nil
		}

		tok = parser.peek()
		if tok.Type == ast.T_COMMA {
			parser.next()
			tok = parser.peek()
		}
	}
	return nil
}

func (parser *Parser) ParseString() ast.Expression {
	var tok = parser.peek()

	if tok.Type == ast.T_QQ_STRING {

		tok = parser.next()
		var str = ast.NewStringWithQuote('"', tok)
		return ast.Expression(str)

	} else if tok.Type == ast.T_Q_STRING {

		tok = parser.next()
		var str = ast.NewStringWithQuote('\'', tok)
		return ast.Expression(str)

	} else if tok.Type == ast.T_IDENT {

		tok = parser.next()
		return ast.Expression(ast.NewString(tok))

	} else if tok.Type == ast.T_INTERPOLATION_START {

		return parser.ParseInterp()

	}
	return nil
}

func (parser *Parser) ParseInterp() ast.Expression {
	debug("ParseInterp at %d", parser.Pos)
	var startTok = parser.peek()

	if startTok.Type != ast.T_INTERPOLATION_START {
		return nil
	}

	parser.accept(ast.T_INTERPOLATION_START)
	var innerExpr = parser.ParseExpression()
	var endTok = parser.expect(ast.T_INTERPOLATION_END)
	var interp = ast.NewInterpolation(innerExpr, startTok, endTok)
	return interp
}

func (parser *Parser) ParseValue() ast.Expression {
	debug("ParseValue")
	var pos = parser.Pos

	// try parse map
	debug("Trying ParseMap")
	if mapValue := parser.ParseMap(); mapValue != nil {
		debug("OK ParseMap")
		return mapValue
	}
	parser.restore(pos)

	debug("ParseList trying")
	if listValue := parser.ParseList(); listValue != nil {
		debug("ParseList OK: %+v", listValue)
		return listValue
	}

	debug("ParseList failed, restoring to %d", pos)
	parser.restore(pos)
	/*
		if stringTerm := parser.ParseInterp(); stringTerm != nil {
			var tok = parser.peek()
			for tok.Type == ast.T_LITERAL_CONCAT {
				var rightExpr = parser.ParseExpression()
				stringTerm = ast.NewBinaryExpression(ast.OpConcat, stringTerm, rightExpr)
				tok = parser.peek()
			}
			return stringTerm
		} else {
			// for other possible string concat expression
		}
	*/
	debug("ParseExpression trying", pos)
	return parser.ParseExpression()
	/*
		var tok = parser.peek()

		// list or map starts with '('
		if tok.Type == ast.T_PAREN_START {
			if expr := parser.ParseExpression(); expr != nil {
				return expr
			}
		}

		tok = parser.peek()
		if tok.Type == ast.T_COLON {
			// it's a map
		}

		tok = parser.peek()
		if tok.Type == ast.T_PAREN_START {
			// parser.ParseMapOrList()
		} else {
			parser.ParseList()
		}
		return nil
	*/
}

func (parser *Parser) ParseList() *ast.List {
	debug("ParseList at %d", parser.Pos)
	var pos = parser.Pos
	var list = parser.ParseCommaSepList()
	if list == nil {
		debug("ParseList failed")
		parser.restore(pos)
		return nil
	}
	return list
}

func (parser *Parser) ParseCommaSepList() *ast.List {
	debug("ParseCommaSepList at %d", parser.Pos)
	var list = ast.NewList()
	list.Separator = ", "

	var tok = parser.peek()
	for tok.Type != ast.T_COMMA && tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_END {

		// when the syntax start with a '(', it could be a list or map.
		if tok.Type == ast.T_PAREN_START {

			parser.next()
			if sublist := parser.ParseCommaSepList(); sublist != nil {
				debug("Appending sublist %+v", list)
				list.Append(sublist)
			}
			// allow empty list here
			parser.expect(ast.T_PAREN_END)

		} else {
			var sublist = parser.ParseSpaceSepList()
			if sublist != nil {
				debug("Appending sublist %+v", list)
				list.Append(sublist)
			} else {
				if list.Len() > 0 {
					return list
				}
				return nil
			}
		}

		if parser.accept(ast.T_COMMA) == nil {
			debug("Returning comma-separated list: %+v\n", list)
			return list
		}
		tok = parser.peek()
	}
	debug("Returning comma-separated list: %+v\n", list)
	// XXX: if there is only one item in the list, we can reduce it to element.
	return list
}

func (parser *Parser) ParseVariable() *ast.Variable {
	var pos = parser.Pos
	var tok = parser.next()
	if tok.Type != ast.T_VARIABLE {
		parser.restore(pos)
		return nil
	}
	return ast.NewVariable(tok)
}

func (parser *Parser) ParseVariableAssignment() ast.Statement {
	var pos = parser.Pos

	var variable = parser.ParseVariable()
	if variable == nil {
		parser.restore(pos)
		return nil
	}

	// skip ":", T_COLON token
	if parser.accept(ast.T_COLON) == nil {
		panic("Expecting colon after variable name")
	}

	var expr = parser.ParseValue()
	if expr == nil {
		panic("Expecting value after variable assignment.")
	}

	parser.expect(ast.T_SEMICOLON)

	// Reduce list or map here
	return ast.NewVariableAssignment(variable, expr)
}

func (parser *Parser) ParseSpaceSepList() *ast.List {
	debug("ParseSpaceSepList at %d", parser.Pos)

	var list = ast.NewList()
	list.Separator = " "

	var tok = parser.peek()

	if tok.Type == ast.T_PAREN_START {
		parser.next()
		if sublist := parser.ParseCommaSepList(); sublist != nil {
			list.Append(sublist)
		}
		parser.expect(ast.T_PAREN_END)
	}

	tok = parser.peek()
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_END {
		var subexpr = parser.ParseExpression()
		if subexpr != nil {
			debug("Parsed Expression: %+v", subexpr)
			list.Append(subexpr)
		} else {
			break
		}
		tok = parser.peek()
		if tok.Type == ast.T_COMMA {
			break
		}
	}
	debug("Returning space-sep list: %+v", list)
	if list.Len() > 0 {
		return list
	}
	return nil
}

/**
We treat the property value section as a list value, which is separated by ',' or ' '
*/
func (parser *Parser) ParsePropertyValue(parentRuleSet *ast.RuleSet, property *ast.Property) *ast.List {
	debug("ParsePropertyValue")
	// var tok = parser.peek()
	var list = ast.NewList()

	var tok = parser.peek()
	for tok.Type != ast.T_SEMICOLON && tok.Type != ast.T_BRACE_END {
		var sublist = parser.ParseList()
		if sublist != nil {
			list.Append(sublist)
			debug("ParsePropertyValue list: %+v", list)
		} else {
			break
		}
		tok = parser.peek()
	}

	tok = parser.peek()
	if tok.Type == ast.T_SEMICOLON || tok.Type == ast.T_BRACE_END {
		parser.next()
	} else {
		panic(fmt.Errorf("Unexpected end of property value. Got %s", tok))
	}
	return list
}

func (parser *Parser) ParseDeclarationBlock(parentRuleSet *ast.RuleSet) *ast.DeclarationBlock {
	var declBlock = ast.DeclarationBlock{}

	var tok = parser.next() // should be '{'
	if tok.Type != ast.T_BRACE_START {
		panic(ParserError{"{", tok.Str})
	}

	tok = parser.next()
	for tok != nil && tok.Type != ast.T_BRACE_END {

		if tok.Type == ast.T_PROPERTY_NAME_TOKEN {
			parser.expect(ast.T_COLON)

			var property = ast.NewProperty(tok)
			var valueList = parser.ParsePropertyValue(parentRuleSet, property)
			_ = valueList
			// property.Values = valueList
			declBlock.Append(property)
			_ = property

		} else if tok.IsSelector() {
			// parse subrule
			panic("subselector unimplemented")
		} else {
			panic("unexpected token")
		}

		tok = parser.next()
	}

	return &declBlock
}

func (parser *Parser) ParseImportStatement() ast.Statement {
	// skip the ast.T_IMPORT token
	var tok = parser.next()

	// Create the import statement node
	var rule = ast.ImportStatement{}

	tok = parser.peek()
	// expecting url(..)
	if tok.Type == ast.T_IDENT {
		parser.advance()

		if tok.Str != "url" {
			panic("invalid function for @import rule.")
		}

		if tok = parser.next(); tok.Type != ast.T_PAREN_START {
			panic("expecting parenthesis after url")
		}

		tok = parser.next()
		rule.Url = ast.Url(tok.Str)

		if tok = parser.next(); tok.Type != ast.T_PAREN_END {
			panic("expecting parenthesis after url")
		}

	} else if tok.IsString() {
		parser.advance()
		rule.Url = ast.RelativeUrl(tok.Str)
	}

	/*
		TODO: parse media query for something like:

		@import url(color.css) screen and (color);
		@import url('landscape.css') screen and (orientation:landscape);
		@import url("bluish.css") projection, tv;
	*/
	tok = parser.peek()
	if tok.Type == ast.T_MEDIA {
		parser.advance()
		rule.MediaList = append(rule.MediaList, tok.Str)
	}

	// must be ast.T_SEMICOLON
	tok = parser.next()
	if tok.Type != ast.T_SEMICOLON {
		panic(ParserError{";", tok.Str})
	}
	return &rule
}
