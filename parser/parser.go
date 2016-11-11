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

import (
	"errors"

	. "github.com/davidli2010/calc/lexer"
)

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

	if token.Kind() == TokenSubOP {
		minus = true
	} else {
		parser.backToken(token)
	}

	if token, err = parser.nextToken(); err != nil {
		return 0, err
	}

	if token.Kind() == TokenNumber {
		value = token.Value()
	} else if token.Kind() == TokenLP {
		if value, err = parser.parseExpression(); err != nil {
			return 0, err
		}

		if token, err = parser.nextToken(); err != nil {
			return 0, err
		}

		if token.Kind() != TokenRP {
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
	var v1 float64
	var err error

	if v1, err = parser.parsePrimaryExpression(); err != nil {
		return 0, err
	}

	for {
		var token *Token

		if token, err = parser.nextToken(); err != nil {
			return 0, err
		}

		if token.Kind() != TokenMulOP && token.Kind() != TokenDivOP {
			parser.backToken(token)
			break
		}

		op := token.Kind()

		var v2 float64
		if v2, err = parser.parsePrimaryExpression(); err != nil {
			return 0, err
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
	var v1 float64
	var err error

	if v1, err = parser.parseTerm(); err != nil {
		return 0, err
	}

	for {
		var token *Token

		if token, err = parser.nextToken(); err != nil {
			return 0, err
		}

		if token.Kind() != TokenAddOP && token.Kind() != TokenSubOP {
			parser.backToken(token)
			break
		}

		op := token.Kind()

		var v2 float64
		if v2, err = parser.parseTerm(); err != nil {
			return 0, err
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
