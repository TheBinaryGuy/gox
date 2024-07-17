package main

import (
	"log"
	"os"
	"time"

	"github.com/TheBinaryGuy/gox/tokenizer"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "Gox (Lox in Go) Interpreter",
		Usage:                "A simple interpreter for the Lox language written in Go",
		Version:              "0.0.1",
		EnableBashCompletion: true,
		Compiled:             time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "TheBinaryGuy",
				Email: "hello@thebinary.dev",
			},
		},

		Commands: []*cli.Command{
			{
				Name:  "tokenize",
				Usage: "tokenize the provided file",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:      "file",
						Value:     "",
						Usage:     "file to tokenize",
						Required:  true,
						TakesFile: true,
						Aliases:   []string{"f"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					file := cCtx.Path("file")
					fileContents, err := os.ReadFile(file)
					if err != nil {
						log.Fatalf("Error reading file: %v\n", err)
					}

					tokens, err := tokenizer.Tokenize(fileContents)
					tokenizer.PrintTokens(tokens)
					if err != nil {
						os.Exit(65)
					}

					return nil
				},
			},
		},
	}

	app.Suggest = true

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
