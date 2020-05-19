package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/minhdanh/vault-template/pkg/template"
)

type vaultConfig struct {
	token, endpoint string
}

type rendererConfig struct {
	vault vaultConfig
}

type renderer struct {
	vaultRenderer *template.VaultTemplateRenderer
}

func NewRenderer(cfg rendererConfig) (*renderer, error) {
	var vaultRenderer *template.VaultTemplateRenderer

	if cfg.vault.token != "" && cfg.vault.endpoint != "" {
		var err error
		vaultRenderer, err = template.NewVaultTemplateRenderer(cfg.vault.token, cfg.vault.endpoint)

		if err != nil {
			return nil, err
		}
	} else {
		panic("Error: Vault endpoint or token is incorrect.")
	}

	return &renderer{
		vaultRenderer: vaultRenderer,
	}, nil
}

func (r *renderer) renderSingleFile(inputFilePath, outputFilePath string) (err error) {
	inputContent, err := ioutil.ReadFile(inputFilePath)

	if err != nil {
		return
	}

	renderedContent := string(inputContent)
	if r.vaultRenderer != nil {
		renderedContent, err = r.vaultRenderer.RenderTemplate(string(inputContent))

		if err != nil {
			return
		}
	}

	// make output path
	outputDirectory := filepath.Dir(outputFilePath)
	err = os.MkdirAll(outputDirectory, 0755)

	if err != nil {
		return
	}
	outputFile, err := os.Create(outputFilePath)

	if err != nil {
		return
	}

	defer func() {
		err = outputFile.Close()
	}()

	_, err = outputFile.Write([]byte(renderedContent))
	fmt.Println("Created file " + outputFilePath)
	return
}

func (r *renderer) renderFromStdinToStdout() (err error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return
	}

	var stdinLines []string
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			stdinLines = append(stdinLines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			if err != nil {
				return err
			}
		}
	}

	renderedContent := strings.Join(stdinLines, "\n")

	if r.vaultRenderer != nil {
		renderedContent, err = r.vaultRenderer.RenderTemplate(renderedContent)

		if err != nil {
			return
		}
	}

	fmt.Printf("%v\n", renderedContent)
	return
}
