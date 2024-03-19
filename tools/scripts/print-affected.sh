#!/bin/bash

set -e

COMMIT_RANGE="origin/$BASE_REF"
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

PROJECT=$(pnpm nx show projects --affected --base="$COMMIT_RANGE" --head="$CURRENT_BRANCH" --select=projects)

echo "$PROJECT"
