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
	"bytes"
	"fmt"
	"strconv"
)

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isSpace(b byte) bool {
	switch b {
	case '\t', '\n', '\v', '\f', '\r', ' ':
		return true
	}
	return false
}

type TokenKind int8

const (
	TokenInvalid = TokenKind(iota)
	TokenNumber
	TokenAddOP
	TokenSubOP
	TokenMulOP
	TokenDivOP
	TokenLP
	TokenRP
	TokenEOL
	TokenEND
)

func (k TokenKind) String() string {
	switch k {
	case TokenInvalid:
		return "Invalid"
	case TokenNumber:
		return "Number"
	case TokenAddOP:
		return "AddOP"
	case TokenSubOP:
		return "SubOP"
	case TokenMulOP:
		return "MulOP"
	case TokenDivOP:
		return "DivOP"
	case TokenLP:
		return "LP"
	case TokenRP:
		return "RP"
	case TokenEOL:
		return "EOL"
	case TokenEND:
		return "END"
	default:
		return "Unknown"
	}
}

type Token struct {
	kind  TokenKind
	str   []byte
	value float64
}

func (token *Token) Kind() TokenKind {
	return token.kind
}

func (token *Token) Value() float64 {
	return token.value
}

func (token *Token) Reset() {
	token.kind = TokenInvalid
	token.str = token.str[:0]
	token.value = 0
}

func (token *Token) Equals(other *Token) bool {
	if token.kind != other.kind {
		return false
	}
	if token.value != other.value {
		return false
	}
	if len(token.str) != len(other.str) {
		return false
	}
	if string(token.str) != string(other.str) {
		return false
	}
	return true
}

func (token *Token) String() string {
	var buf bytes.Buffer
	buf.WriteString("Token{")
	buf.WriteString("kind: ")
	buf.WriteString(token.kind.String())
	buf.WriteString(", str: \"")
	buf.WriteString(string(token.str))
	buf.WriteString("\", value: ")
	buf.WriteString(fmt.Sprint(token.value))
	buf.WriteString("}")
	return buf.String()
}

type lexerStatus int8

const (
	lexInit = lexerStatus(iota)
	lexInt
	lexDot
	lexFrac
)

type Lexer struct {
	line  []byte
	pos   int
	token *Token
}

func NewLexer() *Lexer {
	return &Lexer{token: &Token{}}
}

func (lexer *Lexer) SetLine(line string) {
	lexer.line = []byte(line)
	lexer.pos = 0
}

func (lexer *Lexer) NextToken() (*Token, error) {
	status := lexInit

	token := lexer.token
	token.Reset()

	for lexer.pos < len(lexer.line) {
		char := lexer.line[lexer.pos]

		switch status {
		case lexInt:
			fallthrough
		case lexFrac:
			if char == '.' {
				break
			}
			fallthrough
		case lexDot:
			if !isDigit(char) {
				token.kind = TokenNumber
				value, err := strconv.ParseFloat(string(token.str), 64)
				if err != nil {
					return nil, err
				}
				token.value = value
				return token, nil
			}
		}

		if isSpace(char) {
			lexer.pos++
			if char == '\n' {
				token.kind = TokenEOL
				return token, nil
			}
			continue
		}

		token.str = append(token.str, lexer.line[lexer.pos])
		lexer.pos++

		switch {
		case char == '+':
			token.kind = TokenAddOP
		case char == '-':
			token.kind = TokenSubOP
		case char == '*':
			token.kind = TokenMulOP
		case char == '/':
			token.kind = TokenDivOP
		case char == '(':
			token.kind = TokenLP
		case char == ')':
			token.kind = TokenRP
		case char == '.':
			if status == lexInt {
				status = lexDot
				continue
			} else {
				return nil, fmt.Errorf("syntax error")
			}
		case isDigit(char):
			if status == lexInit {
				status = lexInt
			} else if status == lexDot {
				status = lexFrac
			}
			continue
		default:
			return nil, fmt.Errorf("syntax error")
		}

		return token, nil
	}

	switch status {
	case lexInt:
		fallthrough
	case lexFrac:
		fallthrough
	case lexDot:
		token.kind = TokenNumber
		value, err := strconv.ParseFloat(string(token.str), 64)
		if err != nil {
			return nil, err
		}
		token.value = value
		return token, nil
	default:
		token.kind = TokenEND
		return token, nil
	}
}
