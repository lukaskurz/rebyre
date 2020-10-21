package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

type Clause struct {
	literals []*Literal
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

type Literal struct {
	variable string
	negated  bool
}

func (l *Literal) Negate() {
	l.negated = !l.negated
}

func (l *Literal) Matches(other *Literal) bool {
	return l.variable == other.variable && l.negated == other.negated
}

func (l *Literal) Opposes(other *Literal) bool {
	return l.variable == other.variable && l.negated != other.negated
}

func main() {
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
		clauses[i] = &Clause{literals: parseLiterals(s)}
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
