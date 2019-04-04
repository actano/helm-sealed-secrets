#!/bin/bash

set -ueo pipefail

VERSION=0.1.3

function isAlreadyInstalled() {
  hash sealed-secret-template 2>/dev/null && [[ $(sealed-secret-template -v | cut -d " " -f 3) == ${VERSION} ]]
}

if isAlreadyInstalled; then
  echo "sealed-secret-template is already installed"
else
  echo "Downloading sealed-secret-template version ${VERSION}"
  OS=$(uname | tr '[:upper:]' '[:lower:]')
  URL=https://github.com/actano/helm-sealed-secrets/releases/download/${VERSION}/sealed-secret-template_${OS}_amd64

  curl -sL ${URL} > /usr/local/bin/sealed-secret-template
  chmod +x /usr/local/bin/sealed-secret-template
fi
