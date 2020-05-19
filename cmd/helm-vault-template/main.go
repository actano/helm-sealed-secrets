package main

import (
	"fmt"
	"os"

	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

var Version = "v0.1.3"

func createRenderer(c *cli.Context) (*renderer, error) {
	cfg := rendererConfig{
		vault: vaultConfig{
			endpoint: c.GlobalString("vault.address"),
			token:    c.GlobalString("vault.token"),
		},
	}

	return NewRenderer(cfg)
}

func NewYamlSourceFromFile(file string) func(context *cli.Context) (altsrc.InputSourceContext, error) {
	return func(context *cli.Context) (altsrc.InputSourceContext, error) {
		return altsrc.NewYamlSourceFromFile(file)
	}
}

func main() {
	app := cli.NewApp()
	app.Version = Version
	app.Usage = "Render template with values from Vault to a file"

	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "vault.token",
			Usage:  "Vault API token",
			EnvVar: "VAULT_TOKEN",
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "vault.address",
			Usage:  "Vault API endpoint",
			EnvVar: "VAULT_ADDR",
		}),
	}

	app.Flags = flags

	app.Commands = []cli.Command{
		{
			Name:  "render",
			Usage: "render a template with values from Vault into a file",
			Action: func(c *cli.Context) error {
				renderer, err := createRenderer(c)

				if err != nil {
					return err
				}

				if c.NArg() == 0 {
					stat, err := os.Stdin.Stat()
					if err != nil {
						return err
					}
					if stat.Mode()&os.ModeNamedPipe == 0 {
						cli.ShowCommandHelpAndExit(c, "render", 1)
					} else {
						return renderer.renderFromStdinToStdout()
					}

				} else if c.NArg() == 2 {
					inputFile := c.Args().Get(0)
					outputFile := c.Args().Get(1)

					return renderer.renderSingleFile(inputFile, outputFile)
				} else {
					cli.ShowCommandHelpAndExit(c, "render", 1)
				}
				return nil
			},
			ArgsUsage: "[input file] [output file]",
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering file: %v", err)
		os.Exit(1)
	}
}
