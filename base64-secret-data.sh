#!/bin/sh

set -e

OWN_PATH="${0%/*}"
INPUT_FILE="${OWN_PATH}/secret.yaml"

trap 'rm -r ./tmp' EXIT

mkdir -p ./tmp
yaml2json ${INPUT_FILE} tmp/secret.json

cat ./tmp/secret.json | jq '.data | map_values(. | @base64)' > ./tmp/base64_data.json
cat ./tmp/secret.json | jq --slurpfile base64_data ./tmp/base64_data.json '.data = ($base64_data | .[])' > tmp/secret_base64.json

# output goes to stdout
json2yaml tmp/secret_base64.json
