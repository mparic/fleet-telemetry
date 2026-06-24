#!/bin/sh
set -e

mkdir -p /etc/ssl/fleet-telemetry
printf '%s' "$TLS_SERVER_CERT" | tr -d ' \t\n\r' | base64 -d > /etc/ssl/fleet-telemetry/server.crt
printf '%s' "$TLS_SERVER_KEY"  | tr -d ' \t\n\r' | base64 -d > /etc/ssl/fleet-telemetry/server.key

exec /fleet-telemetry -config /etc/fleet-telemetry/config.json