#!/bin/bash
set -e

# Check precondition for script execution
if [[ ! $(command -v kubeseal) ]]; then
    >&2 echo "Error: kubeseal is not installed"
    exit 1
fi

export OWN_PATH="${0%/*}"
export SCRIPT_BASE64_SECRET_FILE="${OWN_PATH}/base64-secret-data.sh"
DEFAULT_PATHS="${OWN_PATH}/paths.txt"

# Either use the provided paths arguments or the default path file
if [[ -z "$1" ]]
  then
    SECRET_PATHS="${OWN_PATH}/../helm"
    echo "No paths supplied, using the default paths '${SECRET_PATHS}'"
else
    SECRET_PATHS="$@"
fi

function render_template() {
  # Vault settings
  VAULT_ADDR="https://vault.actano.de"
  VAULT_TOKEN_PATH="${HOME}/.vault-token"

  local TEMPLATE_FILE=$1
  local BASENAME_TEMPLATE_FILE=$(basename ${TEMPLATE_FILE})
  local DIRNAME_TEMPLATE_FILE=$(dirname ${TEMPLATE_FILE})

  local NAME="${BASENAME_TEMPLATE_FILE%.template.yaml}"

  # Path of the rendered secret, ignored in git
  SECRET_BASEPATH="tmp_secret"
  SECRET_PATH="${SECRET_BASEPATH}/${NAME}.yaml"
  SECRET_TMP_PATH="${SECRET_BASEPATH}/${NAME}_tmp.yaml"

  # Remove tmp secret folder after script execution
  trap 'rm -r ${SECRET_BASEPATH}' EXIT

  # Render secret template
  docker run --rm \
    -v ${VAULT_TOKEN_PATH}:/root/token \
    -v $(pwd)/${TEMPLATE_FILE}:/root/secret.template.yaml \
    -v $(pwd)/${SECRET_BASEPATH}:/root/${SECRET_BASEPATH} \
    rplan/vault-template:1.0.1 \
      --vault ${VAULT_ADDR} \
      --vault-token-file /root/token \
      --template /root/secret.template.yaml \
      --output /root/${SECRET_PATH}

  # Base64 convert values in data
  # kubeseal doesn't support `stringData`, so we have to convert to base64 beforehand
  docker run --rm \
    -v $(pwd)/${SECRET_PATH}:/data/secret.yaml \
    -v $(pwd)/${SCRIPT_BASE64_SECRET_FILE}:/data/run.sh \
    rplan/transform-yaml \
   > ${SECRET_TMP_PATH}

  # Encrypt secret file with kubeseal
  cat "${SECRET_TMP_PATH}" | kubeseal --controller-namespace sealed-secrets \
  --cert ${OWN_PATH}/cert-rplan-production.pem \
  --format yaml \
  > ${DIRNAME_TEMPLATE_FILE}/../templates/${NAME}.sealed.yaml

  echo "Sealed secret written to ${DIRNAME_TEMPLATE_FILE}/../templates/${NAME}.sealed.yaml"
}

function update_sealed_secrets_in_path() {
  local SECRET_PATH=$1

  # Template files to read vault secret paths from
  export -f render_template
  find ${SECRET_PATH} -name "*.template.yaml" -exec bash -c 'render_template "$0"' {} \;
}

# Update sealed secrets for each path
for SECRET_PATH in ${SECRET_PATHS}
do
  update_sealed_secrets_in_path ${SECRET_PATH}
done
