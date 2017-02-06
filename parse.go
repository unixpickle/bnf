package bnf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"

	"github.com/unixpickle/essentials"
)

// ReadGrammar reads an entire Grammar.
//
// Since this inherently reads until EOF, io.EOF is never
// returned as an error.
//
// If the read fails, the parsed grammar up to that point
// is returned.
//
// The grammar's token references are not verified.
// This means that a symbol might be referenced by never
// defined.
func ReadGrammar(r io.Reader) (g Grammar, err error) {
	defer func() {
		err = essentials.AddCtx("read BNF grammar", err)
	}()
	rr := bufio.NewReader(r)
	line := 1
	for {
		rule, numLines, err := ReadRule(rr)
		if err != nil {
			if err == io.EOF {
				return g, nil
			}
			return g, fmt.Errorf("line %d: %s", line, err.Error())
		}
		line += numLines
		if rule != nil {
			g = append(g, rule)
		}
	}
}

// ReadRule reads the next rule in BNF form.
//
// This returns a nil Rule along with a nil error if the
// next line was blank.
// This returns a nil Rule with an io.EOF error if EOF was
// reached before any non-whitespace characters.
func ReadRule(rr io.RuneReader) (rule *Rule, numLines int, err error) {
	defer func() {
		if err != io.EOF {
			err = essentials.AddCtx("read BNF rule", err)
		}
	}()
	var r reader
	for {
		ch, _, err := rr.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, numLines, err
		}
		if ch == '\n' {
			numLines++
			if r.state == readingRulePreSpace {
				return nil, numLines, nil
			}
		}
		done, err := r.Next(ch)
		if err != nil {
			return nil, numLines, err
		}
		if done {
			return r.Rule(), numLines, nil
		}
	}
	if r.state == readingRulePreSpace {
		return nil, numLines, io.EOF
	}
	numLines++
	if done, err := r.Next('\n'); err != nil {
		return nil, numLines, err
	} else if !done {
		return nil, numLines, errors.New("unexpected EOF")
	} else {
		return r.Rule(), numLines, nil
	}
}

type readerState int

const (
	readingRulePreSpace readerState = iota
	readingRuleName
	readingRuleNameSpace
	readingColon
	readingEqual
	readingExprSpace
	readingExprToken
	readingExprStrNoBackslash
	readingExprStrBackslash
)

type reader struct {
	state readerState

	ruleName        string
	exprs           [][]*Token
	currentExpr     []*Token
	currentExprData string
}

func (r *reader) Rule() *Rule {
	return &Rule{
		Name:        r.ruleName,
		Expressions: r.exprs,
	}
}

func (r *reader) Next(ch rune) (done bool, err error) {
	switch r.state {
	case readingRulePreSpace:
		if unicode.IsSpace(ch) {
			return
		}
		r.ruleName += string(ch)
		r.state = readingRuleName
	case readingRuleName:
		if unicode.IsSpace(ch) {
			r.state = readingRuleNameSpace
		} else {
			r.ruleName += string(ch)
		}
	case readingRuleNameSpace:
		if !unicode.IsSpace(ch) {
			if ch != ':' {
				return true, fmt.Errorf("expected ':' but got %q", ch)
			}
			r.state = readingColon
		}
	case readingColon:
		if ch != ':' {
			return true, fmt.Errorf("expected ':' but got %q", ch)
		}
		r.state = readingEqual
	case readingEqual:
		if ch != '=' {
			return true, fmt.Errorf("expected '=' but got %q", ch)
		}
		r.state = readingExprSpace
	case readingExprSpace:
		if ch == '\n' {
			if len(r.currentExpr) != 0 {
				r.exprs = append(r.exprs, r.currentExpr)
				return true, nil
			}
		} else if ch == '|' {
			if len(r.currentExpr) == 0 {
				return true, errors.New("empty option for rule")
			}
			r.exprs = append(r.exprs, r.currentExpr)
			r.currentExpr = nil
		} else if ch == '"' {
			r.state = readingExprStrNoBackslash
		} else if ch == '<' {
			r.state = readingExprToken
		} else if !unicode.IsSpace(ch) {
			return true, fmt.Errorf("unexpected %q", ch)
		}
	case readingExprToken:
		if ch == '>' {
			r.currentExpr = append(r.currentExpr, &Token{Rule: r.currentExprData})
			r.currentExprData = ""
			r.state = readingExprSpace
		} else {
			r.currentExprData += string(ch)
		}
	case readingExprStrNoBackslash:
		if ch == '"' {
			esc := `"` + r.currentExprData + `"`
			unesc, err := strconv.Unquote(esc)
			if err != nil {
				return true, err
			}
			r.currentExpr = append(r.currentExpr, &Token{Raw: unesc})
			r.currentExprData = ""
			r.state = readingExprSpace
		} else {
			r.currentExprData += string(ch)
			if ch == '\\' {
				r.state = readingExprStrBackslash
			}
		}
	case readingExprStrBackslash:
		r.currentExprData += string(ch)
		r.state = readingExprStrNoBackslash
	default:
		panic("unknown state")
	}
	return false, nil
}
