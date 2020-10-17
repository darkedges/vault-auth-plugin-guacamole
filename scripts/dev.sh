#!/usr/bin/env bash
set -eEuo pipefail

# Setup scratch
SCRATCH="$(pwd)/tmp"
mkdir -p "${SCRATCH}/plugins"

# Build plugin
go build -o "${SCRATCH}/plugins/vault-auth-example"

# Run vault
vault server \
  -dev \
  -dev-plugin-init \
  -dev-plugin-dir "${SCRATCH}/plugins" \
  -dev-root-token-id "root" \
  -log-level "debug" \
  &
sleep 2
VAULT_PID=$!

# Cleanup
function cleanup {
  echo ""
  echo "==> Cleaning up"
  kill -INT "${VAULT_PID}"
  rm -rf "${SCRATCH}"
}
trap cleanup EXIT

# Login
vault login root

vault plugin list
vault auth enable -path=example -plugin-name=vault-auth-example plugin

# Configure
# vault write auth/example/config \
#   username="guacadmin" \
#   password="guacadmin" \
#   url="http://localhost:1234"

# Display config
# vault read auth/example/config

vault write auth/example/login password=super-secret-password

# Wait
wait $!