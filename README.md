# helm-vault-template

This plugin intergrates with Vault to render a template file with Vault paths to a file. Useful to prepare a `values.yaml` file to deploy a helm chart.

### Installation

```bash
helm plugin install https://github.com/cxagroup/helm-vault-template
```

### Usage
```
NAME:
   helm-vault-template - Render a template file

USAGE:
   helm-vault-template [global options] command [arguments...]

VERSION:
   0.0.1

COMMANDS:
     render       Render a template file to a file
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --vault.token value                     Vault token [$VAULT_TOKEN]
   --vault.address value                   Vault API endpoint [$VAULT_ADDR]
   --help, -h                              show help
   --version, -v                           print the version
```

### Examples

Prepare a template file `values.vault.yaml`.

```yaml
config:
  username: {{ vault "secret/data/myservice/admin-user" "username" }}
  password: {{ vault "secret/data/myservice/admin-user" "password" }}
```

To render the template:

```bash
helm vault-template render values.vault.yaml values.yaml
```

`values.yaml` will look like:
```yaml
config:
  username: admin
  password: P@ssw0rd
```

You can also print the output file to stdout:

```bash
helm vault-template render values.vault.yaml -
```
