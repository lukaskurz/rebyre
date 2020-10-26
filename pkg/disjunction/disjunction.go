package disjunction

import (
	"math"
	"regexp"

	"github.com/lukaskurz/rebyre/pkg/literal"
)

// Disjunction to contain disjunction of literals in SAT
type Disjunction struct {
	id       int
	literals []*literal.Literal
	SourceA  int
	SourceB  int
}

// ID returns the id of this disjunction
func (d *Disjunction) ID() int {
	return d.id
}

// Length outputs the length or the "order" of the disjunction
func (d *Disjunction) Length() int {
	return len(d.literals)
}

// IsEmpty checks wether this disjunction is empty i.e. has not literals
func (d *Disjunction) IsEmpty() bool {
	return len(d.literals) == 0
}

// String stringifies the disjunction.
//
// Example: "(!a | b | c)"
func (d *Disjunction) String() string {
	text := "( "

	length := len(d.literals)
	for i, l := range d.literals {
		if l.Negated() {
			text += "!"
		}
		text += l.Variable()
		if i < length-1 {
			text += " | "
		}
	}

	return text + " )"
}

// CompatibleWith checks wether this conjunction is compatible with another in terms of the resolution process.
// It does this by searching for exactly 1 opposing literal and as many as possible matching literals
func (d *Disjunction) CompatibleWith(other *Disjunction) bool {
	matches := 0
	opposed := 0

	for _, l1 := range d.literals {
		for _, l2 := range other.literals {
			if l1.Equals(l2) {
				matches++
				break
			}
			if l1.Opposes(l2) {
				opposed++
				break
			}
		}
	}

	minLength := int(math.Min(float64(d.Length()), float64(other.Length())))
	return opposed == 1 && matches >= minLength-2
}

// Derive derives a disjunction by applying the absorption rule
func (d *Disjunction) Derive(other *Disjunction) *Disjunction {
	var base *Disjunction
	var target *Disjunction

	if d.Length() >= other.Length() {
		base = d
		target = other
	} else {
		base = other
		target = d
	}

	derivation := &Disjunction{id: getNextID(), literals: make([]*literal.Literal, 0), SourceA: base.id, SourceB: target.id}

	var opposer *literal.Literal
	for _, dl := range base.literals {
		for _, ol := range target.literals {
			if dl.Opposes(ol) {
				opposer = dl
				break
			}
		}
		if opposer != nil {
			break
		}
	}

	for _, l := range append(base.literals, target.literals...) {
		if !(l.Equals(opposer) || l.Opposes(opposer)) {
			derivation.literals = append(derivation.literals, l)
		}
	}

	derivation.Sanitize()

	return derivation
}

// Sanitize removes duplicate literals i.e.
//
// Example: "a or b or a" => "a or b"
func (d *Disjunction) Sanitize() {
	clean := false
	for clean == false {
		clean = true

		for i1, l1 := range d.literals {
			for i2, l2 := range d.literals {
				if i1 == i2 {
					continue
				}
				if l1.Equals(l2) {
					clean = false

					// remove element at i2
					d.literals[i2] = d.literals[len(d.literals)-1]
					d.literals = d.literals[:len(d.literals)-1]

					break
				}
			}
			if clean == false {
				break
			}
		}
	}

}

// Equals checks if it is equal to another disjunction, by equaling all literals
func (d *Disjunction) Equals(other *Disjunction) bool {
	if d.Length() != other.Length() {
		return false
	}
	found := false
	for _, l1 := range d.literals {
		found = false
		for _, l2 := range other.literals {
			if l1.Equals(l2) {
				found = true
				break
			}
		}

		// matching literal not found
		if !found {
			return false
		}
	}
	return found
}

// DisjunctionFromString parses a disjunction and the enclosed literals from a string
// Disjunction has to be written in this way:
//
// (a | !!b | !c)
func DisjunctionFromString(text string) (*Disjunction, error) {
	r, err := regexp.Compile("[\\s()]") // match all whitespaces and brackets
	if err != nil {
		return nil, err
	}

	text = r.ReplaceAllString(text, "")
	rLit, err := regexp.Compile(literal.LiteralMatchExp)
	if err != nil {
		return nil, err
	}

	matches := rLit.FindAllString(text, -1)
	literals := make([]*literal.Literal, len(matches))
	for i, m := range matches {
		literals[i], err = literal.LiteralFromString(m)
		if err != nil {
			return nil, err
		}
	}

	return &Disjunction{
		id:       getNextID(),
		literals: literals,
	}, nil
}
