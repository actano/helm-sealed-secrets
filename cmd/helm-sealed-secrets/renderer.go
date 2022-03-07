package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/actano/vault-template/pkg/template"
	"github.com/bmatcuk/doublestar"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type sealedSecretsConfig struct {
	controllerNamespace, publicKey string
}

type vaultConfig struct {
	tokenFile, endpoint string
}

type rendererConfig struct {
	sealedSecrets sealedSecretsConfig
	vault         vaultConfig
}

type renderer struct {
	vaultRenderer *template.VaultTemplateRenderer
	sealedSecrets sealedSecretsConfig
}

var alreadyPrinted = make(map[string]bool)

// printOnce prints any given message only once per program execution
func printOnce(msg string) {
	if alreadyPrinted[msg] != true {
		fmt.Println(msg)
		alreadyPrinted[msg] = true
	}
}

func createErrorWithCause(msg string, cause error) error {
	return fmt.Errorf("%s\n\tcaused by: %s", msg, cause.Error())
}

func NewRenderer(cfg rendererConfig) (*renderer, error) {
	var vaultRenderer *template.VaultTemplateRenderer

	if cfg.vault.tokenFile != "" && cfg.vault.endpoint != "" {
		expandedTokenFile, err := homedir.Expand(cfg.vault.tokenFile)
		if err != nil {
			return nil, err
		}
		vaultToken, err := ioutil.ReadFile(expandedTokenFile)

		if err != nil {
			err = createErrorWithCause("Couldn't read vault token from file: "+expandedTokenFile, err)
			return nil, err
		}

		vaultRenderer, err = template.NewVaultTemplateRenderer(string(vaultToken), cfg.vault.endpoint)

		if err != nil {
			err = createErrorWithCause("Couldn't create vault renderer", err)
			return nil, err
		}
	}

	return &renderer{
		vaultRenderer: vaultRenderer,
		sealedSecrets: cfg.sealedSecrets,
	}, nil
}

func (r *renderer) renderSingleFile(inputFilePath, outputFilePath string) (err error) {
	inputContent, err := ioutil.ReadFile(inputFilePath)

	if err != nil {
		err = createErrorWithCause("Couldn't read input file: "+inputFilePath, err)
		return
	}

	renderedContent := string(inputContent)
	if r.vaultRenderer != nil {
		renderedContent, err = r.vaultRenderer.RenderTemplate(string(inputContent))

		if err != nil {
			err = createErrorWithCause("Couldn't render content from input file: "+inputFilePath, err)
			return
		}
	} else {
		printOnce("NOTE: Not using vault, sealing the secrets as is.")
	}

	base64Data, err := dataToBase64(renderedContent)

	if err != nil {
		err = createErrorWithCause("Couldn't encode rendered content of input file: "+inputFilePath, err)
		return
	}

	sealedContent, err := r.sealSecret(base64Data)

	if err != nil {
		err = createErrorWithCause("Couldn't seal encoded content of input file: "+inputFilePath, err)
		return
	}

	// make output path
	outputDirectory := filepath.Dir(outputFilePath)
	err = os.MkdirAll(outputDirectory, 0755)

	if err != nil {
		err = createErrorWithCause("Couldn't create output directory: "+outputDirectory, err)
		return
	}
	outputFile, err := os.Create(outputFilePath)

	if err != nil {
		err = createErrorWithCause("Couldn't create output file: "+outputFilePath, err)
		return
	}

	defer func() {
		err = outputFile.Close()
	}()

	_, err = outputFile.Write([]byte(sealedContent))
	fmt.Println("Created sealed file " + outputFilePath)
	return
}

func (r *renderer) sealSecret(secret string) (sealedSecret string, err error) {
	args := []string{"--format", "yaml"}
	if r.sealedSecrets.publicKey != "" {
		args = append(args, "--cert", r.sealedSecrets.publicKey)
	} else if r.sealedSecrets.controllerNamespace != "" {
		args = append(args, "--controller-namespace", r.sealedSecrets.controllerNamespace)
	}
	cmd := exec.Command("kubeseal", args...)
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
	if err != nil {
		switch err := err.(type) {
		case *exec.Error:
			fmt.Println(err)
		case *exec.ExitError:
			fmt.Println("kubeseal returned error:", string(err.Stderr))
		}
		return
	}
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
				valueBytes := []byte(fmt.Sprintf("%v", dataItem.Value))
				data[k].Value = base64.StdEncoding.EncodeToString(valueBytes)
			}
		}
	}

	out, err := yaml.Marshal(secret)

	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (r *renderer) renderDir(inputDir, outputDir string) error {
	matches, err := findFiles(inputDir, "*.template.yaml")

	if err != nil {
		return err
	}

	if len(matches) == 0 {
		fmt.Printf("No matching files found in input directory %q", inputDir)
		return nil
	}

	inputOutputPaths, err := GetInputOutputPaths(matches, inputDir, outputDir)

	if err != nil {
		return err
	}

	for _, match := range inputOutputPaths {
		err = r.renderSingleFile(match.InputPath, match.OutputPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetInputOutputPaths(matches []string, inputDir, outputDir string) (inputOutputPaths []InputOutputPaths, err error) {
	for _, match := range matches {
		var relativePath string
		relativePath, err = filepath.Rel(inputDir, match)
		if err != nil {
			return
		}
		subPath := filepath.Dir(relativePath)
		inputFilename := filepath.Base(relativePath)
		outputFileName := strings.Replace(inputFilename, ".template.yaml", ".sealed.yaml", 1)
		outputFilePath := filepath.Join(outputDir, subPath, outputFileName)
		inputOutputPaths = append(inputOutputPaths, InputOutputPaths{
			InputPath:  match,
			OutputPath: outputFilePath,
		})
	}
	return
}

func findFiles(targetDir, pattern string) ([]string, error) {
	globPattern := filepath.Join(targetDir, "**", pattern)
	return doublestar.Glob(globPattern)
}

type InputOutputPaths struct {
	InputPath, OutputPath string
}
