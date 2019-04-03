package main

import (
    "encoding/base64"
    "fmt"
    "github.com/Luzifer/rconfig"
    "github.com/actano/vault-template/pkg/template"
    "gopkg.in/yaml.v2"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
)

type config struct {
    VaultEndpoint  string `flag:"vault-address" env:"VAULT_ADDR" default:"https://127.0.0.1:8200" description:"Vault API endpoint. Also configurable via VAULT_ADDR."`
    VaultTokenFile string `flag:"vault-token-file" env:"VAULT_TOKEN_FILE" description:"The file which contains the vault token. Also configurable via VAULT_TOKEN_FILE."`
    InputFile      string `flag:"input-file,i" description:"The input secret template which should be rendered and sealed."`
    OutputFile     string `flag:"output-file,o" description:"The output file path where the sealed secret should be written to."`
}

func usage(msg string) {
    println(msg)
    rconfig.Usage()
    os.Exit(1)
}

func parseConfig() (*config, error) {
    cfg := &config{}
    err := rconfig.Parse(cfg)

    if err != nil {
        return nil, err
    }

    if cfg.VaultEndpoint == "" {
        return nil, fmt.Errorf("no vault endpoint give")
    }

    return cfg, nil
}

func renderSingleFile(renderer *template.VaultTemplateRenderer, inputFilePath, outputFilePath string) (err error) {
    inputContent, err := ioutil.ReadFile(inputFilePath)

    if err != nil {
        return
    }

    renderedContent, err := renderer.RenderTemplate(string(inputContent))

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

    outputFile, err := os.Create(outputFilePath)

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

func main() {
    config, err := parseConfig()

    if err != nil {
    	usage("")
       panic(err)
    }

    vaultToken, err := ioutil.ReadFile(config.VaultTokenFile)

    if err != nil {
       panic(err)
    }

    renderer, err := template.NewVaultTemplateRenderer(string(vaultToken), "https://vault.actano.de")

    if err != nil {
       panic(err)
    }

    if config.InputFile != "" {
       err = renderSingleFile(renderer, config.InputFile, config.OutputFile)
       if err != nil {
           panic(err)
       }
       return
    }
}
