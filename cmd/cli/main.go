package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/lukaskurz/rebyre/pkg/disjunction"
)

var (
	index int
	out   io.StringWriter
)

func main() {
	index = 0

	solveCommand := &cli.Command{
		Name:    "solve",
		Aliases: []string{"s"},
		Usage:   "rebyre solve <path/to/file.bool>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "output",
				Aliases:   []string{"o"},
				Usage:     "output sets where solutions are streamed. keep it empty for STD (terminal output) or provide a file path",
				Required:  false,
				Hidden:    false,
				TakesFile: true,
			},
		},
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

			oFlag := c.String("output")
			if len(strings.TrimSpace(oFlag)) == 0 {
				out = os.Stdout
			} else {
				abs, err := filepath.Abs(oFlag)
				if err != nil {
					return err
				}

				f, err := os.Create(abs)
				if err != nil {
					return err
				}
				out = f
				defer func() {
					err = f.Close()
					fmt.Println(err)
				}()
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
				out.WriteString(fmt.Sprintf("\nSolution #%d\n\n", i))
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
				Name:  "PhD Lukas G. Kurz",
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
	text := d.String()

	if len(indent) > 0 {
		if left {
			out.WriteString("T")
		} else {
			out.WriteString(fmt.Sprintf("%s└", indent))
		}
	}
	out.WriteString(text)

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
		out.WriteString("\n")
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
		out.WriteString(fmt.Sprintf("%d %s\n", d.ID(), d.String()))
	}
}

func printCombinations(combinations []*disjunction.Disjunction) {
	for _, c := range combinations {
		out.WriteString(fmt.Sprintf("%d %s %d %d\n", c.ID(), c.String(), c.SourceA, c.SourceB))
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
