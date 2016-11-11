// Copyright 2016 David Li
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package parser

import "errors"

type Parser struct {
	lexer           *Lexer
	lookAhead       Token
	lookAheadExists bool
}

func NewParser() *Parser {
	return &Parser{lexer: NewLexer()}
}

func (parser *Parser) reset() {
	parser.lookAheadExists = false
}

func (parser *Parser) nextToken() (*Token, error) {
	if parser.lookAheadExists {
		parser.lookAheadExists = false
		return &parser.lookAhead, nil
	} else {
		return parser.lexer.NextToken()
	}
}

func (parser *Parser) backToken(token *Token) {
	parser.lookAhead = *token
	parser.lookAheadExists = true
}

/*
 * primary_expression
 *     : TokenNumber
 *     | TokenSubOP primary_expression
 *     | TokenLP expression TokenRP
 */
func (parser *Parser) parsePrimaryExpression() (float64, error) {
	var token *Token
	var err error
	var minus bool
	var value float64

	if token, err = parser.nextToken(); err != nil {
		return 0, err
	}

	if token.kind == TokenSubOP {
		minus = true
	} else {
		parser.backToken(token)
	}

	if token, err = parser.nextToken(); err != nil {
		return 0, err
	}

	if token.kind == TokenNumber {
		value = token.value
	} else if token.kind == TokenLP {
		if value, err = parser.parseExpression(); err != nil {
			return 0, err
		}

		if token, err = parser.nextToken(); err != nil {
			return 0, err
		}

		if token.kind != TokenRP {
			return 0, errors.New("missing ')'")
		}
	} else {
		return 0, errors.New("syntax error")
	}

	if minus {
		value = -value
	}

	return value, nil
}

/*
 * term
 *     : primary_expression
 *     | term TokenMulOP primary_expression
 *     | term TokenDivOP primary_expression
 */
func (parser *Parser) parseTerm() (float64, error) {
	v1, err1 := parser.parsePrimaryExpression()
	if err1 != nil {
		return 0, err1
	}

	for {
		token, err := parser.nextToken()
		if err != nil {
			return 0, err
		}

		if token.kind != TokenMulOP && token.kind != TokenDivOP {
			parser.backToken(token)
			break
		}

		op := token.kind

		v2, err2 := parser.parsePrimaryExpression()
		if err2 != nil {
			return 0, nil
		}

		switch op {
		case TokenMulOP:
			v1 *= v2
		case TokenDivOP:
			v1 /= v2
		}
	}

	return v1, nil
}

/*
 * expression
 *     : term
 *     | expression TokenAddOP term
 *     | expression TokenSubOP term
 */
func (parser *Parser) parseExpression() (float64, error) {
	v1, err1 := parser.parseTerm()
	if err1 != nil {
		return 0, err1
	}

	for {
		token, err := parser.nextToken()
		if err != nil {
			return 0, err
		}

		if token.kind != TokenAddOP && token.kind != TokenSubOP {
			parser.backToken(token)
			break
		}

		op := token.kind

		v2, err2 := parser.parseTerm()
		if err2 != nil {
			return 0, err2
		}

		switch op {
		case TokenAddOP:
			v1 += v2
		case TokenSubOP:
			v1 -= v2
		}
	}

	return v1, nil
}

func (parser *Parser) Parse(line string) (float64, error) {
	parser.reset()
	parser.lexer.SetLine(line)
	return parser.parseExpression()
}
