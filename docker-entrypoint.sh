#!/bin/sh
set -e

mkdir -p /etc/ssl/fleet-telemetry
echo "$TLS_SERVER_CERT" | base64 -d > /etc/ssl/fleet-telemetry/server.crt
echo "$TLS_SERVER_KEY"  | base64 -d > /etc/ssl/fleet-telemetry/server.key

exec /fleet-telemetry -config /etc/fleet-telemetry/config.json