package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/lukaskurz/rebyre/pkg/disjunction"
)

var index int

func main() {
	index = 0

	solveCommand := &cli.Command{
		Name:    "solve",
		Aliases: []string{"s"},
		Usage:   "rebyre solve <path/to/file.bool>",
		Action: func(c *cli.Context) error {
			verbose := c.Bool("verbose")
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

			if verbose {
				printDisjunctions(disjunctions)
			}
			fmt.Println("Starting resolution:")

			emptyClauses := make([]*disjunction.Disjunction, 0)
			for len(emptyClauses) == 0 {
				combinations := combineDisjunctions(disjunctions)
				if verbose {
					printCombinations(combinations)
				}

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
		Compiled:             time.Date(2020, time.October, 25, 19, 37, 0, 0, time.UTC),
		Description:          "You either know what this thing does or you don't. Repo is found at https://github.com/lukaskurz/rebyre",
		EnableBashCompletion: true,
		Version:              "4.20.69",
		Commands: []*cli.Command{
			solveCommand,
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "verbose"},
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

func printTree(all []*disjunction.Disjunction, d *disjunction.Disjunction, indent string, left bool) {
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

func getDisjunction(id int, all []*disjunction.Disjunction) *disjunction.Disjunction {
	for _, e := range all {
		if e.ID() == id {
			return e
		}
	}

	return nil
}

func printDisjunctions(disjunctions []*disjunction.Disjunction) {
	for _, d := range disjunctions {
		fmt.Println(fmt.Sprintf("%d %s", d.ID(), d.ToString()))
	}
}

func printCombinations(combinations []*disjunction.Disjunction) {
	for _, c := range combinations {
		fmt.Println(fmt.Sprintf("%d %s %d %d", c.ID(), c.ToString(), c.SourceA, c.SourceB))
	}
}

func parseDisjunctions(text string) ([]*disjunction.Disjunction, error) {
	splitted := strings.Split(text, "&")
	disjunctions := make([]*disjunction.Disjunction, len(splitted))

	for i, s := range splitted {
		var err error
		disjunctions[i], err = disjunction.DisjunctionFromString(s)
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

func getEmptyClauses(disjunctions []*disjunction.Disjunction) []*disjunction.Disjunction {
	clauses := make([]*disjunction.Disjunction, 0)

	for _, d := range disjunctions {
		if d.IsEmpty() {
			clauses = append(clauses, d)
		}
	}
	return clauses
}

func combineDisjunctions(disjunctions []*disjunction.Disjunction) []*disjunction.Disjunction {
	length := len(disjunctions)

	combinations := make([]*disjunction.Disjunction, 0)

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

func isClauseContained(clauses []*disjunction.Disjunction, clause *disjunction.Disjunction) bool {
	for _, c := range clauses {
		if c.Equals(clause) {
			return true
		}
	}
	return false
}
