package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "rebyre",
		Usage: "Tool to do a refutation by resolution on a proposition in CNF",
		Action: func(c *cli.Context) error {
			fmt.Println("asd asd asd asd asd asd asd")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
