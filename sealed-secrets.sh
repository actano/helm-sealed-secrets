#!/bin/bash

function usage() {
  cat <<EOF
Available commands:
  update      Render sealed secrets
EOF
}

case "${1:-help}" in
  --help|-h|help)
    usage
    ;;
  *)
    usage
    exit 1
    ;;
esac
