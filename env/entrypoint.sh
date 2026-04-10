#!/bin/sh

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CONFIG_FILE="$SCRIPT_DIR/env.${MODE:-dev}.json"

echo "MODE=$MODE"
echo "Looking for config: $CONFIG_FILE"

ls -la "$SCRIPT_DIR"

if [ ! -f "$CONFIG_FILE" ]; then
  echo "❌ Config file not found: $CONFIG_FILE"
  exit 1
fi

echo "✅ Config file found"

CONTENT="$(cat "$CONFIG_FILE")"

if [ -z "$CONTENT" ]; then
  echo "❌ Config file is empty"
  exit 1
fi

export APP_CONFIG="$(echo "$CONTENT" | tr -d '\n')"

echo "APP_CONFIG length: ${#APP_CONFIG}"

exec "$@"
