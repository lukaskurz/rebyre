package literal

import (
	"reflect"
	"testing"
)

func TestLiteralEquals(t *testing.T) {
	a := &Literal{
		variable: "a",
		negated:  false,
	}

	sameA := &Literal{
		variable: "a",
		negated:  false,
	}

	t.Run("same equal", func(t *testing.T) {
		if !a.Equals(sameA) {
			t.Error("same equal FAILED, expected true")
		}
	})

	b := &Literal{
		variable: "b",
		negated:  false,
	}

	bNeg := &Literal{
		variable: "b",
		negated:  true,
	}

	aNeg := &Literal{
		variable: "a",
		negated:  true,
	}

	t.Run("not equal", func(t *testing.T) {
		if a.Equals(b) {
			t.Error("not equal FAILED, expected false for a == b")
		}

		if a.Equals(bNeg) {
			t.Error("not equal FAILED, expected false for a == b!")
		}

		if a.Equals(aNeg) {
			t.Error("not equal FAILED, expected false for a == !a")
		}
	})
}

func TestLiteralOpposes(t *testing.T) {
	a := &Literal{
		variable: "a",
		negated:  false,
	}

	sameA := &Literal{
		variable: "a",
		negated:  false,
	}

	b := &Literal{
		variable: "b",
		negated:  false,
	}

	bNeg := &Literal{
		variable: "b",
		negated:  true,
	}

	aNeg := &Literal{
		variable: "a",
		negated:  true,
	}

	t.Run("same variable", func(t *testing.T) {
		if a.Opposes(sameA) {
			t.Error("expected false for a opposes a")
		}

		if !a.Opposes(aNeg) {
			t.Error("expected true for a opposes !a")
		}
	})

	t.Run("different variable", func(t *testing.T) {
		if a.Opposes(b) {
			t.Error("expected false for a opposes b")
		}

		if a.Opposes(bNeg) {
			t.Error("expected false for a opposes !b")
		}
	})
}

func TestLiteralToString(t *testing.T) {
	a := &Literal{
		variable: "a",
		negated:  false,
	}

	b := &Literal{
		variable: "b",
		negated:  false,
	}

	bNeg := &Literal{
		variable: "b",
		negated:  true,
	}

	aNeg := &Literal{
		variable: "a",
		negated:  true,
	}

	toString := []string{
		a.ToString(),
		aNeg.ToString(),
		b.ToString(),
		bNeg.ToString(),
	}

	expected := []string{
		"a",
		"!a",
		"b",
		"!b",
	}

	for i, e := range toString {
		if e != expected[i] {
			t.Errorf("FAILED, expected toString[%d] to be %s, not %s", i, expected[i], e)
		}
	}
}

func TestLiteralFromString(t *testing.T) {
	t.Run("invalid literals", func(t *testing.T) {
		invalids := []string{
			"!",
			"1",
			"1asd",
			"as1",
			"aaa1aa",
			"!1",
			"!a1",
			"!1a",
			"!a1a",
			"!1asdasd1",
			"a!",
			"asfsf!",
		}

		for _, i := range invalids {
			_, err := LiteralFromString(i)
			if err == nil {
				t.Errorf("FAILED, expected error for \"%s\"", i)
			}
		}
	})

	t.Run("valid literals", func(t *testing.T) {
		valids := []struct {
			s string
			l *Literal
		}{
			{"a", &Literal{variable: "a", negated: false}},
			{"mythical", &Literal{variable: "mythical", negated: false}},
			{"!myth", &Literal{variable: "myth", negated: true}},
			{"!!a", &Literal{variable: "a", negated: false}},
			{"!a", &Literal{variable: "a", negated: true}},
		}

		for _, i := range valids {
			lit, err := LiteralFromString(i.s)
			if err != nil {
				t.Errorf("FAILED, expected no error for \"%s\"", i.s)
			}

			if !reflect.DeepEqual(*lit, *i.l) {
				t.Errorf("FAILED, expected same literal for \"%s\"", i.s)
			}
		}
	})
}
