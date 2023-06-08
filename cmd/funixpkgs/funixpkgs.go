package main

import (
	"errors"
	"fmt"
	"funixpkgs/src/nixpkgs"
	"log"
	"os"

	"github.com/urfave/cli"
)

const description = `Find deprecated packages in nixpkgs from GitHub.`

func main() {

	app := cli.NewApp()
	app.Name = description
	app.Usage = "Find which packages are not updated"

	app.Commands = []cli.Command{
		{
			Name:        "about",
			HelpName:    "about",
			Action:      about,
			ArgsUsage:   ` `,
			Usage:       `show information about the application and author.`,
			Description: `show desc.`,
		},
		{
			Name:      "find",
			HelpName:  "find",
			Action:    findUnupdatedPkgs,
			ArgsUsage: ` `,
			Usage:     `find which packages are not updated.`,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "path",
					Usage: "path of nixpkgs.",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func about(c *cli.Context) error {

	fmt.Println("Author: Quentin JOLY")
	fmt.Println("Version: 1.0.0")
	fmt.Println("License: MIT")
	fmt.Println("Description: " + description)
	return nil
}

func findUnupdatedPkgs(c *cli.Context) error {

	if !c.IsSet("path") && c.String("path") == "" {
		return errors.New("path flag must be provided")
	}
	path := c.String("path")

	infoPath, err := os.Stat(path)
	if err != nil || !infoPath.IsDir() {
		if os.IsNotExist(err) {
			return errors.New("path does not exist")
		}
		if !infoPath.IsDir() {
			return errors.New("path is not a directory")
		}
	}

	listPkgs, err := nixpkgs.FindUnupdatedPkgs(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(listPkgs))

	return nil
}
