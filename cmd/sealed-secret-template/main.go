package main

import (
    "github.com/Luzifer/rconfig"
    "os"
)

type config struct {
    VaultEndpoint                    string `flag:"vault-address" env:"VAULT_ADDR" default:"https://127.0.0.1:8200" description:"Vault API endpoint. Also configurable via VAULT_ADDR."`
    VaultTokenFile                   string `flag:"vault-token-file" env:"VAULT_TOKEN_FILE" description:"The file which contains the vault token. Also configurable via VAULT_TOKEN_FILE."`
    SealedSecretsControllerNamespace string `flag:",c" description:"Sealed secret controller namespace"`
    InputFile                        string `flag:"input-file,i" description:"The input secret template which should be rendered and sealed."`
    OutputFile                       string `flag:"output-file,o" description:"The output file path where the sealed secret should be written to."`
}

func usage(msg string) {
    println(msg)
    rconfig.Usage()
    os.Exit(1)
}

func parseConfig() (*config, error) {
    cfg := &config{}
    err := rconfig.Parse(cfg)
    return cfg, err
}

func main() {
    cfg, err := parseConfig()

    if err != nil {
    	usage("")
        panic(err)
    }

    renderer, err := NewRenderer(cfg)

    if err != nil {
		panic(err)
    }

    if cfg.InputFile != "" && cfg.OutputFile != "" {
        err = renderer.renderSingleFile()
        if err != nil {
            panic(err)
        }
        return
    }
}
