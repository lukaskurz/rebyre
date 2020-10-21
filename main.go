package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var counter int

type Clause struct {
	id       int
	literals []*Literal
}

func (c *Clause) Length() int {
	return len(c.literals)
}

func (c *Clause) IsEmpty() bool {
	return len(c.literals) == 0
}

func (c *Clause) ToString() string {
	text := "( "

	length := len(c.literals)
	for i, l := range c.literals {
		if l.negated {
			text += "!"
		}
		text += l.variable
		if i < length-1 {
			text += " | "
		}
	}

	return text + " )"
}

func (c *Clause) CompatibleWith(other *Clause) bool {
	matches := 0
	opposed := 0

	for _, l1 := range c.literals {
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

	minLength := int(math.Min(float64(c.Length()), float64(other.Length())))
	return opposed == 1 && matches >= minLength-1
}

func (c *Clause) Derive(other *Clause) *Clause {
	// TODO:
	return &Clause{literals: make([]*Literal, 0)}
}

func (c *Clause) Equals(other *Clause) bool {
	if c.Length() != other.Length() {
		return false
	}
	found := false
	for _, l1 := range c.literals {
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

type Literal struct {
	variable string
	negated  bool
}

func (l *Literal) Equals(other *Literal) bool {
	return l.variable == other.variable && l.negated == other.negated
}

func (l *Literal) Opposes(other *Literal) bool {
	return l.variable == other.variable && l.negated != other.negated
}

func main() {
	counter = 0

	app := &cli.App{
		Name: "rebyre",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "stdout",
				Usage:   "output result to file",
			},
		},
		Usage: "Tool to do a refutation by resolution on a proposition in CNF",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("No file input specified")
			}
			text, err := readTextFromFile(c.Args().First())
			if err != nil {
				return err
			}

			clauses := parseClauses(text)
			for _, c := range clauses {
				fmt.Println(c.ToString())
			}

			for !containsEmptyClauses(clauses) {
				clauses = combineClauses(clauses)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func parseClauses(text string) []*Clause {
	splitted := strings.Split(text, "&")
	clauses := make([]*Clause, len(splitted))

	for i, s := range splitted {
		clauses[i] = &Clause{
			id:       getNextId(),
			literals: parseLiterals(s),
		}
	}

	return clauses
}

func parseLiterals(text string) []*Literal {
	text = strings.ReplaceAll(text, "(", "")
	text = strings.ReplaceAll(text, ")", "")
	splitted := strings.Split(text, "|")
	literals := make([]*Literal, len(splitted))

	for i, s := range splitted {
		literals[i] = parseLiteral(s)
	}

	return literals
}

func parseLiteral(text string) *Literal {
	if len(text) == 2 {
		return &Literal{variable: strings.ReplaceAll(text, "!", ""), negated: true}
	}
	return &Literal{variable: text, negated: false}
}

func readTextFromFile(filepath string) (string, error) {
	buffer, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	text := string(buffer)
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\r", "")
	text = strings.ReplaceAll(text, " ", "")

	return text, nil
}

func getNextId() int {
	counter++
	return counter
}

func containsEmptyClauses(clauses []*Clause) bool {
	for _, c := range clauses {
		if c.IsEmpty() {
			return true
		}
	}
	return false
}

func combineClauses(clauses []*Clause) []*Clause {
	for _, c1 := range clauses {
		for _, c2 := range clauses {
			if c1.CompatibleWith(c2) {
				derived := c1.Derive(c2)
				if !isClauseContained(clauses, derived) {
					derived.id = getNextId()
					clauses = append(clauses, derived)
				}
			}
		}
	}

	return clauses
}

func isClauseContained(clauses []*Clause, clause *Clause) bool {
	for _, c := range clauses {
		if c.Equals(clause) {
			return true
		}
	}
	return false
}
