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

import "testing"

func TestParser(t *testing.T) {
	lines := [...]struct {
		line  string
		value float64
	}{
		{"1+2", 3},
		{"1 - 2\n", -1},
		{"2*3\n", 6},
		{" 10 / 5", 2},
		{"- ( -1 - 2 )", 3},
		{"1.5 + 2. * (-1 - 2.2) / 10", 0.86},
	}

	parser := NewParser()

	for _, expr := range lines {
		if value, err := parser.Parse(expr.line); err != nil {
			t.Errorf("error: %v, line: %s, expect value: %f",
				err, expr.line, expr.value)
		} else if value != expr.value {
			t.Errorf("line: %s, expect value: %f, actual value: %f",
				expr.line, expr.value, value)
		}
	}
}
