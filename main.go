package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var index int

func main() {
	counter = 0
	index = 0

	solveCommand := &cli.Command{
		Name:  "solve",
		Usage: "rebyre solve <path/to/file.bool>",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("No file input specified")
			}
			text, err := readTextFromFile(c.Args().First())
			if err != nil {
				return err
			}

			disjunctions, err := parseDisjunctions(text)
			if err != nil {
				return err
			}

			// printDisjunctions(disjunctions)
			fmt.Println("Starting")

			emptyClauses := make([]*Disjunction, 0)
			for len(emptyClauses) == 0 {
				combinations := combineDisjunctions(disjunctions)
				// printCombinations(combinations)

				disjunctions = append(disjunctions, combinations...)

				emptyClauses = getEmptyClauses(disjunctions)
			}

			fmt.Println("Found an empty clause !!")

			for i, e := range emptyClauses {
				fmt.Printf("\nSolution #%d\n\n", i)
				printTree(disjunctions, e, "", true)
			}

			return nil
		},
	}

	app := &cli.App{
		Name:                 "rebyre",
		Description:          "You either knnow what this thing does or you don't. Repo is found at https://github.com/lukaskurz/rebyre",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			solveCommand,
		},
		Authors: []*cli.Author{
			{
				Name:  "Lukas G. Kurz",
				Email: "me@lukaskurz.com",
			},
		},
		Usage: "Tool to do a refutation by resolution on a proposition in CNF. Input has to be in CNF i.e. ( a | b | !c ) & ( !a | b)",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func printTree(all []*Disjunction, d *Disjunction, indent string, left bool) {
	text := d.ToString()

	if len(indent) > 0 {
		if left {
			fmt.Print("┬")
		} else {
			fmt.Printf("%s└", indent)
		}
	}
	fmt.Print(text)

	nextIndent := indent
	for i := 0; i < len(text); i++ {
		nextIndent += " "
	}

	if !left {
		nextIndent += " "
	}

	if d.SourceA != 0 {
		next := getDisjunction(d.SourceA, all)
		printTree(all, next, nextIndent+"|", true)
	}
	if d.SourceB != 0 {
		next := getDisjunction(d.SourceB, all)
		printTree(all, next, nextIndent, false)
	}
	if d.SourceA == 0 && d.SourceB == 0 {
		fmt.Println()
	}

}

func getDisjunction(id int, all []*Disjunction) *Disjunction {
	for _, e := range all {
		if e.id == id {
			return e
		}
	}

	return nil
}

func printDisjunctions(disjunctions []*Disjunction) {
	for _, d := range disjunctions {
		fmt.Println(fmt.Sprintf("%d %s", d.id, d.ToString()))
	}
}

func printCombinations(combinations []*Disjunction) {
	for _, c := range combinations {
		fmt.Println(fmt.Sprintf("%d %s %d %d", c.id, c.ToString(), c.SourceA, c.SourceB))
	}
}

func parseDisjunctions(text string) ([]*Disjunction, error) {
	splitted := strings.Split(text, "&")
	disjunctions := make([]*Disjunction, len(splitted))

	for i, s := range splitted {
		var err error
		disjunctions[i], err = DisjunctionFromString(s)
		if err != nil {
			return nil, err
		}
	}

	return disjunctions, nil
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

func getEmptyClauses(disjunctions []*Disjunction) []*Disjunction {
	clauses := make([]*Disjunction, 0)

	for _, d := range disjunctions {
		if d.IsEmpty() {
			clauses = append(clauses, d)
		}
	}
	return clauses
}

func combineDisjunctions(disjunctions []*Disjunction) []*Disjunction {
	length := len(disjunctions)

	combinations := make([]*Disjunction, 0)

	for _, base := range disjunctions[index:] {
		for _, target := range disjunctions {
			if base.CompatibleWith(target) {
				derived := base.Derive(target)
				if !isClauseContained(disjunctions, derived) && !isClauseContained(combinations, derived) {
					combinations = append(combinations, derived)
				}
			}
		}
	}

	index = length

	return combinations
}

func isClauseContained(clauses []*Disjunction, clause *Disjunction) bool {
	for _, c := range clauses {
		if c.Equals(clause) {
			return true
		}
	}
	return false
}
