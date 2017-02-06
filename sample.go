package bnf

import (
	"fmt"
	"math/rand"
	"strings"
)

// Sample creates a string by sampling randomly from the
// grammar.
func Sample(g Grammar, rootRule string) (string, error) {
	r := g.Get(rootRule)
	if r == nil {
		return "", fmt.Errorf("sample grammar: rule %q not found", rootRule)
	}
	if len(r.Expressions) == 0 {
		return "", fmt.Errorf("sample grammar: no expressions for %q", r.Expressions)
	}
	e := r.Expressions[rand.Intn(len(r.Expressions))]
	var parts []string
	for _, t := range e {
		if t.Rule == "" {
			parts = append(parts, t.Raw)
		} else {
			part, err := Sample(g, t.Rule)
			if err != nil {
				return "", err
			}
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, ""), nil
}
