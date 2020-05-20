#!/bin/bash

set -ueo pipefail

VERSION=v0.1.3

export PLUGIN_DIR=$(dirname "$0")

function isAlreadyInstalled() {
  [ -f ${PLUGIN_DIR}/helm-vault-template ] && [[ $(${PLUGIN_DIR}/helm-vault-template -v | cut -d " " -f 3) == ${VERSION} ]]
}

if isAlreadyInstalled; then
  echo "helm-vault-template is already installed"
else
  echo "Downloading helm-vault-template version ${VERSION}"
  OS=$(uname | tr '[:upper:]' '[:lower:]')
  URL=https://github.com/cxagroup/helm-vault-template/releases/download/${VERSION}/helm-vault-template_${OS}_amd64

  temp_file=$(mktemp)
  trap "rm ${temp_file}" EXIT

  statuscode=$(curl -w "%{http_code}" -L ${URL} -o ${temp_file})

  if [[ ! "${statuscode}" == "200" ]]; then
    echo "Failed to download binary"
    exit 1
  fi

  cp ${temp_file} ${PLUGIN_DIR}/helm-vault-template
  chmod +x ${PLUGIN_DIR}/helm-vault-template
fi
