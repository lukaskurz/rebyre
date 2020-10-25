package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

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

func parseClauses(text string) []*Disjunction {
	splitted := strings.Split(text, "&")
	clauses := make([]*Disjunction, len(splitted))

	for i, s := range splitted {
		clauses[i] = &Disjunction{
			id:       getNextID(),
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

func containsEmptyClauses(clauses []*Disjunction) bool {
	for _, c := range clauses {
		if c.IsEmpty() {
			return true
		}
	}
	return false
}

func combineClauses(clauses []*Disjunction) []*Disjunction {
	for _, c1 := range clauses {
		for _, c2 := range clauses {
			if c1.CompatibleWith(c2) {
				derived := c1.Derive(c2)
				if !isClauseContained(clauses, derived) {
					derived.id = getNextID()
					clauses = append(clauses, derived)
				}
			}
		}
	}

	return clauses
}

func isClauseContained(clauses []*Disjunction, clause *Disjunction) bool {
	for _, c := range clauses {
		if c.Equals(clause) {
			return true
		}
	}
	return false
}
