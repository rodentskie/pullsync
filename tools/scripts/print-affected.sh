#!/bin/bash

set -e

COMMIT_RANGE="origin/$BASE_REF"
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
EXCLUDE="${EXCLUDE:-cybersec.app,cybersec.app-e2e}"

PROJECT=$(./node_modules/.bin/nx show projects --affected --base="$COMMIT_RANGE" --head="$CURRENT_BRANCH" --select=projects --exclude "$EXCLUDE")

echo "$PROJECT"
