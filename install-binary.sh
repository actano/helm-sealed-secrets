#!/bin/bash

set -ueo pipefail

VERSION=0.1.1

if hash sealed-secret-template 2>/dev/null; then
  echo "sealed-secret-template is already installed"
else
  OS=$(uname | tr '[:upper:]' '[:lower:]')
  URL=https://github.com/actano/helm-sealed-secrets/releases/download/${VERSION}/sealed-secret-template_${OS}_amd64

  curl -sL ${URL} > /usr/local/bin/sealed-secret-template
  chmod +x /usr/local/bin/sealed-secret-template
fi
