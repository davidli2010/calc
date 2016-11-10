/*
Copyright 2016 David Li

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package calc

import "testing"

func TestLexer(t *testing.T) {
	line := "1.5 + 2. * (-1 - 2.2) / 10\n"
	tokens := [...]Token{
		Token{TokenNumber, []byte("1.5"), float64(1.5)},
		Token{TokenAddOP, []byte("+"), 0},
		Token{TokenNumber, []byte("2."), float64(2)},
		Token{TokenMulOP, []byte("*"), 0},
		Token{TokenLP, []byte("("), 0},
		Token{TokenSubOP, []byte("-"), 0},
		Token{TokenNumber, []byte("1"), float64(1)},
		Token{TokenSubOP, []byte("-"), 0},
		Token{TokenNumber, []byte("2.2"), float64(2.2)},
		Token{TokenRP, []byte(")"), 0},
		Token{TokenDivOP, []byte("/"), 0},
		Token{TokenNumber, []byte("10"), float64(10)},
		Token{TokenEOL, []byte(""), 0},
	}

	lexer := NewLexer()
	lexer.SetLine(line)
	for i := 0; ; i++ {
		if token, err := lexer.NextToken(); err != nil {
			t.Fatal(err)
		} else {
			if !token.Equals(&tokens[i]) {
				t.Fatalf("expect token [%v], actual token [%v]\n", token, &tokens[i])
			}

			if token.kind == TokenEOL {
				break
			}
		}
	}
}
