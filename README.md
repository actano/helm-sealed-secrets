# helm-sealed-secrets

This plugin is used to generate sealed secrets out of secret. It supports template files with vault paths.
This way, you can store both the template and their rendered representation in git.

### Installation

```bash
helm install https://github.com/actano/helm-sealed-secrets
```

### Usage
| Parameter                                    | Description                          | Default                                                                                        |
| ---------------------------------            | ------------------------------------ | ----------------------------------------------------------------------------                   |
| `-c, --sealed-secrets-controller-namespace`  | Sealed secrets controller namespace  | `kube-system` as defined in [`sealed-secrets`](https://github.com/bitnami-labs/sealed-secrets) |
| `--vault-address`, env `VAULT_ADDR`          | Vault endpoint address               | `https://127.0.0.1:8200`                                                                       |
| `--vault-token-file`, env `VAULT_TOKEN_FILE` | Your personal vault token file       | Not set                                                                                        |


### Examples

Read these examples to see how the program works.

#### Single file

Specify a secret template `my-secret.template.yaml`.
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-secret-name
  namespace: dev
type: Opaque
data:
  username: {{ vault "secret/myservice/admin-user" "username" }}
  password: {{ vault "secret/myservice/admin-user" "password" }}
```

Executing

```bash
helm sealed-secrets -i my-secret.template.yaml -o my-secret.yaml
```

gives you a file `my-secret.yaml`

```yaml
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  creationTimestamp: null
  name: my-secret-name
  namespace: dev
spec:
  encryptedData:
    username: 7tgrVWorKLqoZc...
    password: LbeaMTWxTpWAKD...
```
