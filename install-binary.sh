#!/bin/bash

set -ueo pipefail

VERSION=0.16.0
KUBESEAL_VERSION=0.16.0

if ! hash kubeseal 2>/dev/null; then
  echo "kubeseal not found. Installing kubeseal"
  if [[ "$(uname)" == "Darwin" ]]; then
    brew install kubeseal
  elif [[ "$(uname)" == "Linux" ]]; then
    temp_file=$(mktemp)
    trap "rm ${temp_file}" EXIT
    statuscode=$(curl -w "%{http_code}" -sL "https://github.com/bitnami-labs/sealed-secrets/releases/download/v${KUBESEAL_VERSION}/kubeseal-linux-amd64" -o ${temp_file})

    if [[ ! "${statuscode}" == "200" ]]; then
      echo "Failed to download kubeseal: Status code ${statuscode}"
      exit 1
    fi

    cp ${temp_file} /usr/local/bin/kubeseal
    chmod +x /usr/local/bin/kubeseal
  fi
fi

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

  cp ${temp_file} /usr/local/bin/helm-sealed-secrets
  chmod +x /usr/local/bin/helm-sealed-secrets
fi
