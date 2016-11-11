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

package lexer

import "testing"

func TestLexer(t *testing.T) {
	line := "1.5 + 2. * (-1 - 2.2) / 10 \n"
	tokens := [...]Token{
		{TokenNumber, []byte("1.5"), float64(1.5)},
		{TokenAddOP, []byte("+"), 0},
		{TokenNumber, []byte("2."), float64(2)},
		{TokenMulOP, []byte("*"), 0},
		{TokenLP, []byte("("), 0},
		{TokenSubOP, []byte("-"), 0},
		{TokenNumber, []byte("1"), float64(1)},
		{TokenSubOP, []byte("-"), 0},
		{TokenNumber, []byte("2.2"), float64(2.2)},
		{TokenRP, []byte(")"), 0},
		{TokenDivOP, []byte("/"), 0},
		{TokenNumber, []byte("10"), float64(10)},
		{TokenEOL, []byte(""), 0},
		{TokenEND, []byte(""), 0},
	}

	lexer := NewLexer()
	lexer.SetLine(line)
	for _, tk := range tokens {
		if token, err := lexer.NextToken(); err != nil {
			t.Fatal(err)
		} else {
			if !token.Equals(&tk) {
				t.Fatalf("expect token [%v], actual token [%v]\n", &tk, token)
			}

			if token.kind == TokenEND {
				break
			}
		}
	}
}
