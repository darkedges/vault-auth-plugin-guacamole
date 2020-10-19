#!/usr/bin/env bash
set -eEuo pipefail

# Setup scratch
SCRATCH="$(pwd)/tmp"
mkdir -p "${SCRATCH}/plugins"

# Build plugin
go build -o "${SCRATCH}/plugins/vault-auth-guacamole"

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
vault auth enable -path=guacamole -plugin-name=vault-auth-guacamole plugin

# Display info
vault path-help auth/guacamole/
vault read auth/guacamole/info

# Configure
vault write auth/guacamole/config \
  access_token="guacadmin" \
  teams="guacadmin"

# Display config
vault read auth/guacamole/config

vault write auth/guacamole/login password=super-secret-password

# Wait
wait $!