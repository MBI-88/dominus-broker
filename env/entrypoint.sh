#!/bin/sh

export APP_CONFIG="$(cat env.${MODE:-dev}.json | tr -d '\n')"

exec "$@"
