#!/bin/bash

set -e

COMMIT_RANGE="origin/$BASE_REF"
AFFECTED="${TARGET:-apps}"
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

./node_modules/.bin/nx affected -t "$AFFECTED" --base="$COMMIT_RANGE" --head="$CURRENT_BRANCH" --parallel=6 --output-style=stream
