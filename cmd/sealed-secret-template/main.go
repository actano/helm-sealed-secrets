package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

func createRenderer(c *cli.Context) (*renderer, error) {
	vaultEndpoint := c.GlobalString("vault.address")
	vaultTokenFile := c.GlobalString("vault.token-file")
	sealedSecretsControllerNamespace := c.GlobalString("sealed-secrets.controller-namespace")

	return NewRenderer(sealedSecretsControllerNamespace, vaultTokenFile, vaultEndpoint)
}

func NewYamlSourceFromFile(file string) func(context *cli.Context) (altsrc.InputSourceContext, error) {
	return func(context *cli.Context) (altsrc.InputSourceContext, error) {
		return altsrc.NewYamlSourceFromFile(file)
	}
}

func MaybeLoadConfigFromFile(flags []cli.Flag) cli.BeforeFunc {
	return func(c *cli.Context) error {
		configFilePath := c.String("config-file")

		if _, err := os.Stat(configFilePath); err == nil {
			return altsrc.InitInputSourceWithContext(flags, NewYamlSourceFromFile(configFilePath))(c)
		}

		return nil
	}
}

func main() {
	app := cli.NewApp()
	app.HideVersion = true

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "config-file",
			Usage: "Config file to configure the other flags",
			Value: ".sealed-secrets.yaml",
		},
		cli.StringFlag{
			Name:  "vault.token-file",
			Usage: "Location of the vault token file",
			Value: "~/.vault-token",
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "vault.address",
			Usage:  "Vault API endpoint",
			EnvVar: "VAULT_ADDR",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "sealed-secrets.controller-namespace",
			Usage: "The namespace in which the sealed secrets controller runs",
		}),
	}

	app.Before = MaybeLoadConfigFromFile(flags)
	app.Flags = flags

	app.Commands = []cli.Command{
		{
			Name:  "enc",
			Usage: "encrypt a secret template into a sealed secret",
			Action: func(c *cli.Context) error {
				if c.NArg() < 2 {
					cli.ShowCommandHelpAndExit(c, "enc", 1)
				}

				renderer, err := createRenderer(c)

				if err != nil {
					return err
				}

				inputFile := c.Args().Get(0)
				outputFile := c.Args().Get(1)

				return renderer.renderSingleFile(inputFile, outputFile)
			},
			ArgsUsage: "[input file] [output file]",
		},
		{
			Name:      "enc-dir",
			ArgsUsage: "[input directory] [output directory]",
			Usage:     "encrypt all secret templates in a directory structure",
			UsageText: "Encrypts all files with the pattern '*.template.yaml' in the given input directory including all subdirectories. The sealed secrets will be written to the given output directory according to the same directory structure as the input directory.",
			Action: func(c *cli.Context) error {
				if c.NArg() < 2 {
					cli.ShowCommandHelpAndExit(c, "enc-dir", 1)
				}

				renderer, err := createRenderer(c)

				if err != nil {
					return err
				}

				inputDir := c.Args().Get(0)
				outputDir := c.Args().Get(1)

				return renderer.renderDir(inputDir, outputDir)
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		panic(err)
	}
}
