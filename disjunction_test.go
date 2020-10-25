package main

import (
	"testing"
)

func setup() []*Disjunction {
	a := &Literal{variable: "a", negated: false}
	notA := &Literal{variable: "a", negated: true}
	b := &Literal{variable: "b", negated: false}
	notB := &Literal{variable: "b", negated: true}
	c := &Literal{variable: "c", negated: false}
	notC := &Literal{variable: "c", negated: true}

	disjunctions := []*Disjunction{
		{id: 0, literals: []*Literal{
			a,
			notB,
			c,
		}},
		{id: 1, literals: []*Literal{
			c,
		}},
		{id: 2, literals: []*Literal{
			notA,
			notC,
			b,
		}},
		{id: 3, literals: []*Literal{}},
	}

	return disjunctions
}

func setup1() []*Disjunction {
	a := &Literal{variable: "a", negated: false}
	notA := &Literal{variable: "a", negated: true}
	b := &Literal{variable: "b", negated: false}
	notB := &Literal{variable: "b", negated: true}
	c := &Literal{variable: "c", negated: false}
	notC := &Literal{variable: "c", negated: true}

	disjunctions := []*Disjunction{
		{id: 0, literals: []*Literal{
			a,
			notB,
			c,
		}},
		{id: 1, literals: []*Literal{
			a,
			b,
			c,
		}},
		{id: 2, literals: []*Literal{
			notA,
			notB,
			c,
		}},
		{id: 3, literals: []*Literal{
			a,
			b,
			notC,
		}},
		{id: 4, literals: []*Literal{
			notA,
			b,
			c,
		}},
	}

	return disjunctions
}

func setup2() (sources []*Disjunction, derivations []*Disjunction) {
	a := &Literal{variable: "a", negated: false}
	notA := &Literal{variable: "a", negated: true}
	b := &Literal{variable: "b", negated: false}
	notB := &Literal{variable: "b", negated: true}
	c := &Literal{variable: "c", negated: false}

	sources = []*Disjunction{
		{id: 0, literals: []*Literal{
			a,
			notB,
			c,
		}},
		{id: 1, literals: []*Literal{
			a,
			b,
			c,
		}},
		{id: 2, literals: []*Literal{
			notA,
			notB,
			c,
		}},
		{id: 3, literals: []*Literal{
			a,
			b,
		}},
		{id: 4, literals: []*Literal{
			notA,
			b,
		}},
	}

	derivations = []*Disjunction{
		{literals: []*Literal{a, c}},
		{literals: []*Literal{notB, c}},
		{literals: []*Literal{b}},
	}

	return sources, derivations

}

func TestDisjunctionLength(t *testing.T) {
	disjunctions := setup()

	lengths := []int{
		disjunctions[0].Length(),
		disjunctions[1].Length(),
		disjunctions[2].Length(),
		disjunctions[3].Length(),
	}

	results := []int{
		3,
		1,
		3,
		0,
	}

	for i, e := range lengths {
		if e != results[i] {
			t.Errorf("FAILED, expected length of disjunction[%d] to be %d, not %d", i, results[i], lengths[i])
		}
	}
}

func TestDisjunctionIsEmpty(t *testing.T) {
	disjunctions := setup()

	emptiness := []bool{
		disjunctions[0].IsEmpty(),
		disjunctions[1].IsEmpty(),
		disjunctions[2].IsEmpty(),
		disjunctions[3].IsEmpty(),
	}

	results := []bool{
		false,
		false,
		false,
		true,
	}

	for i, e := range emptiness {
		if e != results[i] {
			t.Errorf("FAILED, expected emptiness of disjunction[%d] to be %t, not %t", i, results[i], emptiness[i])
		}
	}
}

func TestDisjunctionToString(t *testing.T) {
	disjunctions := setup()

	toStrings := []string{
		disjunctions[0].ToString(),
		disjunctions[1].ToString(),
		disjunctions[2].ToString(),
		disjunctions[3].ToString(),
	}

	results := []string{
		"( a | !b | c )",
		"( c )",
		"( !a | !c | b )",
		"(  )",
	}

	for i, e := range toStrings {
		if e != results[i] {
			t.Errorf("FAILED, expected string of disjunction[%d] to be %s, not %s", i, results[i], toStrings[i])
		}
	}
}

func TestDisjunctionCompatibleWith(t *testing.T) {
	disjunctions := setup1()

	compatability0 := []bool{
		disjunctions[0].CompatibleWith(disjunctions[0]),
		disjunctions[0].CompatibleWith(disjunctions[1]),
		disjunctions[0].CompatibleWith(disjunctions[2]),
		disjunctions[0].CompatibleWith(disjunctions[3]),
		disjunctions[0].CompatibleWith(disjunctions[4]),
	}

	compatability1 := []bool{
		disjunctions[1].CompatibleWith(disjunctions[0]),
		disjunctions[1].CompatibleWith(disjunctions[1]),
		disjunctions[1].CompatibleWith(disjunctions[2]),
		disjunctions[1].CompatibleWith(disjunctions[3]),
		disjunctions[1].CompatibleWith(disjunctions[4]),
	}

	results0 := []bool{
		false,
		true,
		true,
		false,
		false,
	}

	results1 := []bool{
		true,
		false,
		false,
		true,
		true,
	}

	for i, e := range compatability0 {
		if e != results0[i] {
			t.Errorf("FAILED, expected compatability of disjunction[0]&[%d] to be %t, not %t", i, results0[i], e)
		}
	}

	for i, e := range compatability1 {
		if e != results1[i] {
			t.Errorf("FAILED, expected compatability of disjunction[1]&[%d] to be %t, not %t", i, results1[i], e)
		}
	}
}

func TestDisjunctionDerive(t *testing.T) {
	sources, expected := setup2()

	derivations := []*Disjunction{
		sources[0].Derive(sources[1]),
		sources[0].Derive(sources[2]),
		sources[3].Derive(sources[4]),
	}

	for i, e := range derivations {
		if !expected[i].Equals(e) {
			t.Errorf("FAILED, derivation[%d] is not correct", i)
		}
	}
}

func TestDisjunctionSanitize(t *testing.T) {
	s0, err := DisjunctionFromString("a | a | b")
	s1, err := DisjunctionFromString("b | c | b")
	s2, err := DisjunctionFromString("a | c | b")
	s3, err := DisjunctionFromString("!b | !b | a")
	s4, err := DisjunctionFromString("!b | !b")
	if err != nil {
		t.Errorf("FAILED, got an error: %s", err.Error())
	}

	s0.Sanitize()
	s1.Sanitize()
	s2.Sanitize()
	s3.Sanitize()
	s4.Sanitize()

	sources := []*Disjunction{s0, s1, s2, s3, s4}

	d0, err := DisjunctionFromString("a | b")
	d1, err := DisjunctionFromString("c | b")
	d2, err := DisjunctionFromString("a | c | b")
	d3, err := DisjunctionFromString("!b | a")
	d4, err := DisjunctionFromString("!b ")
	if err != nil {
		t.Errorf("FAILED, got an error: %s", err.Error())
	}

	expected := []*Disjunction{d0, d1, d2, d3, d4}

	for i, e := range sources {
		if !e.Equals(expected[i]) {
			t.Errorf("FAILED, expected sources[%d] and expected[%d] to be equal", i, i)
		}
	}

}

func TestDisjunctionEquals(t *testing.T) {
	d0 := setup()
	d1 := setup1()

	if !d0[0].Equals(d1[0]) {
		t.Errorf("FAILED, expected d0[0] and d1[0] to be equal")
	}
	if d0[1].Equals(d1[1]) {
		t.Errorf("FAILED, expected d0[1] and d1[1] not to be equal")
	}
}

func TestGetNextID(t *testing.T) {
	counter = 0
	id1 := getNextID()
	id2 := getNextID()
	if id1 != 1 {
		t.Errorf("FAILED, expected id1 to be %d not %d", 1, id1)
	}
	if id2 != 2 {
		t.Errorf("FAILED, expected id2 to be %d not %d", 2, id2)
	}
}

func TestDisjunctionFromString(t *testing.T) {
	sources := []string{
		"(a | !b | c)",
		"c",
		"!a |!c| b ",
	}
	targets := setup()

	for i, e := range sources {
		result, err := DisjunctionFromString(e)
		if err != nil {
			t.Errorf("FAILED, got an an error with \"%s\": %s ", e, err.Error())
		}
		if !result.Equals(targets[i]) {
			t.Errorf("FAILED, expected sources[%d]: \"%s\" and targets[%d] to be equal", i, sources[i], i)
		}
	}
}
