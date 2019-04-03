package main

import (
	"encoding/base64"
	"github.com/actano/vault-template/pkg/template"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type renderer struct {
	vaultRenderer *template.VaultTemplateRenderer
	cfg           *config
}

func NewRenderer(cfg *config) (*renderer, error) {
	var vaultRenderer *template.VaultTemplateRenderer

	if cfg.VaultTokenFile != "" && cfg.VaultEndpoint != "" {
		vaultToken, err := ioutil.ReadFile(cfg.VaultTokenFile)

		if err != nil {
			panic(err)
		}

		vaultRenderer, err = template.NewVaultTemplateRenderer(string(vaultToken), "https://vault.actano.de")

		if err != nil {
			return nil, err
		}
	}

	return &renderer{
		vaultRenderer: vaultRenderer,
		cfg:           cfg,
	}, nil
}

func (r *renderer) renderSingleFile() (err error) {
	inputContent, err := ioutil.ReadFile(r.cfg.InputFile)

	if err != nil {
		return
	}

	renderedContent, err := r.vaultRenderer.RenderTemplate(string(inputContent))

	if err != nil {
		return
	}

	base64Data, err := dataToBase64(renderedContent)

	if err != nil {
		return
	}

	sealedContent, err := sealSecret(base64Data)

	if err != nil {
		return
	}

	outputFile, err := os.Create(r.cfg.OutputFile)

	if err != nil {
		return
	}

	defer func() {
		err = outputFile.Close()
	}()

	_, err = outputFile.Write([]byte(sealedContent))
	return
}

func sealSecret(secret string) (sealedSecret string, err error) {
	cmd := exec.Command("kubeseal", "--controller-namespace", "sealed-secrets", "--format", "yaml")
	stdin, err := cmd.StdinPipe()

	if err != nil {
		return
	}

	go func() {
		defer func() {
			err = stdin.Close()
		}()
		_, err = io.WriteString(stdin, secret)
	}()

	out, err := cmd.Output()
	sealedSecret = string(out)

	return
}

func dataToBase64(secretContent string) (string, error) {
	secret := yaml.MapSlice{}
	err := yaml.Unmarshal([]byte(secretContent), &secret)

	if err != nil {
		return "", err
	}

	for _, item := range secret {
		if item.Key == "data" {
			data := item.Value.(yaml.MapSlice)
			for k, dataItem := range data {
				value := dataItem.Value.(string)
				data[k].Value = base64.StdEncoding.EncodeToString([]byte(value))
			}
		}
	}

	out, err := yaml.Marshal(secret)

	if err != nil {
		return "", err
	}

	return string(out), nil
}

