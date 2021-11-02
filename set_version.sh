#!/bin/bash

set -ueo pipefail

if [[ $# < 1 ]]; then
  echo "Usage: $0 <version>"
  exit 1
fi

if [[  $(git diff --stat) != '' ]]; then
  echo 'Please commit your local changes first.'
  exit 1
fi

VERSION=$1

echo "Updating to version $VERSION. Please press Enter to start or CMD+C to cancel"
read ignore

sed -i "s/^VERSION=.*/VERSION=$VERSION/g" install-binary.sh
sed -i "s/version: .*/version: $VERSION/g" plugin.yaml
sed -i "s/VERSION=.*/VERSION=$VERSION/g" Makefile

git commit -am "Bump version to \`$VERSION\`"
git tag $VERSION
