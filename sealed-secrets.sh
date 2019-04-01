#!/bin/bash

HELM_BIN="${HELM_BIN:-helm}"

function usage() {
  cat <<EOF
Available commands:
  update      Render sealed secrets
EOF
}

function is_help() {
  case "$1" in
    --help|-h|help)
      return 0
      ;;
    *)
      return 1
      ;;
  esac
}

function update_usage() {
  cat <<EOF
Update sealed secrets

Example:
  $ ${HELM_BIN} sealed-secrets update <input path>
EOF
}

function update() {
  if is_help "$1"; then
    update_usage
    return
  fi
}

case "${1:-help}" in
  update)
    if [[ $# -lt 1 ]]; then
      update_usage
      echo "Error: input directory required"
      exit 1
    fi
    update "$2"
    ;;
  --help|-h|help)
    usage
    ;;
  *)
    usage
    exit 1
    ;;
esac
