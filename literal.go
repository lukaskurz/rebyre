package main

import (
	"fmt"
	"regexp"
	"strings"
)

// LiteralMatchExp is a regexp to match a single literal
const LiteralMatchExp = "!*[a-zA-Z]+"

// LiteralParseExp is a regexp to match components of a literal
const LiteralParseExp = "([!]*)([a-zA-Z]+)"

// Literal is struct to contain a SAT literal
type Literal struct {
	variable string
	negated  bool
}

// Equals checks if another literal is equal
func (l *Literal) Equals(other *Literal) bool {
	return l.variable == other.variable && l.negated == other.negated
}

// Opposes checks if another literal is the same variable but opposite negation
func (l *Literal) Opposes(other *Literal) bool {
	return l.variable == other.variable && l.negated != other.negated
}

// ToString prints the literal as string
func (l *Literal) ToString() string {
	text := ""
	if l.negated {
		text += "!"
	}
	text += l.variable
	return text
}

// LiteralFromString parses a Literal a string
func LiteralFromString(text string) (*Literal, error) {
	lit := &Literal{}
	text = strings.ReplaceAll(text, " ", "")
	length := len(text)
	if length == 0 {
		return nil, fmt.Errorf("Literal cannot be an empty string")
	}

	// Match a text optionally preceded by negation signs and ending in a variable.
	// Only a-zA-Z for variable name allowed
	//
	// Example: !!a or !a or b
	//
	expr, err := regexp.Compile(fmt.Sprintf("^%s$", LiteralParseExp))
	if err != nil {
		return nil, err
	}

	match := expr.MatchString(text)
	if !match {
		return nil, fmt.Errorf("string is of invalid syntax: %s", text)
	}

	sub := expr.FindStringSubmatch(text)
	noOfNeg := len(sub[1])
	lit.negated = noOfNeg%2 != 0
	lit.variable = sub[2]

	return lit, nil
}
