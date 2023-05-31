package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "template",
				Usage: "Set up template for commit message",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("TODO")
					return nil
				},
			},
			{
				Name:  "clear",
				Usage: "Unset current template",
				Action: func(cCtx *cli.Context) error {
					fmt.Println("TODO")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
