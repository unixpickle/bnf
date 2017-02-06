package bnf

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseSuccess(t *testing.T) {
	var grammar = `
<rule1> ::= <tok1>
<rule2> ::= "hey"
<rule3> ::= <tok1> 'hey'
<rule4> ::= <tok1> | "hey"
<rule5> ::= <tok1> "hey" | 'hey' <tok3>
<YoRule-Hey> ::= <myToken-name> 'hey"quotes' |
	<indented-token> | "yo\" |
	<hey>
	`
	g, err := ReadGrammar(bytes.NewBufferString(grammar))
	if err != nil {
		t.Fatal(err)
	}
	expected := Grammar([]*Rule{
		{"rule1", [][]*Token{{{Rule: "tok1"}}}},
		{"rule2", [][]*Token{{{Raw: "hey"}}}},
		{"rule3", [][]*Token{{{Rule: "tok1"}, {Raw: "hey"}}}},
		{"rule4", [][]*Token{{{Rule: "tok1"}}, {{Raw: "hey"}}}},
		{"rule5", [][]*Token{{{Rule: "tok1"}, {Raw: "hey"}}, {{Raw: "hey"}, {Rule: "tok3"}}}},
		{"YoRule-Hey", [][]*Token{{{Rule: "myToken-name"}, {Raw: "hey\"quotes"}},
			{{Rule: "indented-token"}}, {{Raw: "yo\\"}}, {{Rule: "hey"}}}},
	})
	if len(g) != len(expected) {
		t.Fatalf("expected %d but got %d rules", len(expected), len(g))
	}
	for i, x := range expected {
		a := g[i]
		if !reflect.DeepEqual(a, x) {
			t.Errorf("rule %d: expected %v but got %v", i, x, a)
		}
	}
}
