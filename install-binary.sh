#!/bin/bash

set -ueo pipefail

VERSION=1.17.10

function isAlreadyInstalled() {
  hash helm-sealed-secrets 2>/dev/null && [[ $(helm-sealed-secrets -v | cut -d " " -f 3) == ${VERSION} ]]
}

if isAlreadyInstalled; then
  echo "helm-sealed-secrets is already installed"
else
  echo "Downloading helm-sealed-secrets version ${VERSION}"
  OS=$(uname | tr '[:upper:]' '[:lower:]')
  URL=https://github.com/actano/helm-sealed-secrets/releases/download/${VERSION}/helm-sealed-secrets_${OS}_amd64

  temp_file=$(mktemp)
  trap "rm ${temp_file}" EXIT

  statuscode=$(curl -w "%{http_code}" -sL ${URL} -o ${temp_file})

  if [[ ! "${statuscode}" == "200" ]]; then
    echo "Failed to download binary"
    exit 1
  fi

  cp ${temp_file} $HELM_PLUGIN_DIR/sealed-secrets.sh
  chmod +x $HELM_PLUGIN_DIR/sealed-secrets.sh
fi
