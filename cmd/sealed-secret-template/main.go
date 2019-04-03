package main

import (
	"github.com/Luzifer/rconfig"
	"github.com/bmatcuk/doublestar"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
    VaultEndpoint                    string `flag:"vault-address" env:"VAULT_ADDR" default:"https://127.0.0.1:8200" description:"Optional. Vault API endpoint. Also configurable via VAULT_ADDR."`
    VaultTokenFile                   string `flag:"vault-token-file" env:"VAULT_TOKEN_FILE" description:"Optional. The file which contains the vault token. Also configurable via VAULT_TOKEN_FILE."`
    SealedSecretsControllerNamespace string `flag:"sealed-secrets-controller-namespace,c" description:"Sealed secret controller namespace"`
    InputFile                        string `flag:"input-file,i" description:"The input secret template which should be rendered and sealed."`
    OutputFile                       string `flag:"output-file,o" description:"The output file path where the sealed secret should be written to."`
    InputDir                         string `flag:"input-dir,I" description:"The directory in which to find secret templates to render and seal. Files must match the pattern '<filename>.template.yaml'. The folder structure will be preserved and created at the configured output dir."`
    OutputDir                        string `flag:"output-dir,O" description:"The directory in which to put the rendered sealed secret files. The directory structure from the input directory will be preserved."`
}

func printUsage(msg string) {
    println(msg)
    rconfig.Usage()
    os.Exit(1)
}

func parseConfig() (*config, error) {
    cfg := &config{}
    err := rconfig.Parse(cfg)
    return cfg, err
}

func findFiles(targetDir, pattern string) ([]string, error) {
	globPattern := filepath.Join(targetDir, "**", pattern)
	return doublestar.Glob(globPattern)
}

type InputOutputPaths struct {
	InputPath, OutputPath string
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
		outputFileName := strings.Replace(inputFilename, ".template.yaml", ".yaml", 1)
		outputFilePath := filepath.Join(outputDir, subPath, outputFileName)
		inputOutputPaths = append(inputOutputPaths, InputOutputPaths {
			InputPath:  match,
			OutputPath: outputFilePath,
		})
	}
	return
}

func main() {
    cfg, err := parseConfig()

    if err != nil {
    	printUsage("")
        panic(err)
    }

    renderer, err := NewRenderer(cfg)

    if err != nil {
		panic(err)
    }

    if cfg.InputFile != "" && cfg.OutputFile != "" {
        err = renderer.renderSingleFile(cfg.InputFile, cfg.OutputFile)
        if err != nil {
            panic(err)
        }
        return
    }

    if cfg.InputDir != "" && cfg.OutputDir != "" {
    	matches, err := findFiles(cfg.InputDir, "*.template.yaml")

		if err != nil {
			panic(err)
		}

    	inputOutputPaths, err := GetInputOutputPaths(matches, cfg.InputDir, cfg.OutputDir)

    	if err != nil {
    		panic(err)
		}
    	for _, match := range inputOutputPaths {
    		err = renderer.renderSingleFile(match.InputPath, match.OutputPath)
    		if err != nil {
    			panic(err)
			}
		}
    	return
    }

    printUsage("")
}
