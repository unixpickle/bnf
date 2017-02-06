package bnf

import (
	"fmt"
	"strings"
)

// A Grammar is a set of rules.
type Grammar []*Rule

// String returns the BNF of the grammar.
func (g Grammar) String() string {
	var lines []string
	for _, r := range g {
		lines = append(lines, r.String())
	}
	return strings.Join(lines, "\n")
}

// Get gets a rule by name.
// It returns nil if no rule is found.
func (g Grammar) Get(name string) *Rule {
	for _, x := range g {
		if x.Name == name {
			return x
		}
	}
	return nil
}

// A Rule is a rule for a given token name.
type Rule struct {
	Name string

	// Each element corresponds to an option for the rule.
	Expressions [][]*Token
}

// String returns the BNF of the rule.
func (r *Rule) String() string {
	var exprStrs []string
	for _, e := range r.Expressions {
		var tokStrs []string
		for _, t := range e {
			tokStrs = append(tokStrs, t.String())
		}
		exprStrs = append(exprStrs, strings.Join(tokStrs, " "))
	}
	return r.Name + " ::= " + strings.Join(exprStrs, " | ")
}

// A Token is either a reference to a rule, or a raw
// character sequence.
type Token struct {
	// If non-empty, this token references a rule.
	Rule string

	Raw string
}

// String returns the BNF for the token.
func (t *Token) String() string {
	if t.Rule != "" {
		return "<" + t.Rule + ">"
	} else {
		return fmt.Sprintf("%+q", t.Raw)
	}
}
