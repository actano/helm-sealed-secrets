#!/bin/bash

set -ueo pipefail

if [[ $# < 1 ]]; then
  echo "Usage: $0 <version>"
  exit 1
fi

VERSION=$1

sed -i "" "s/VERSION=.*/VERSION=${VERSION}/g" install-binary.sh
sed -i "" "s/version: .*/version: ${VERSION}/g" plugin.yaml
sed -i "" "s/VERSION=.*/VERSION=${VERSION}/g" Makefile
