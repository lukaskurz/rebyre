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

			fmt.Println(text)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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
