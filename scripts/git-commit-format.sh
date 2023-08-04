#!/bin/bash
#
# Format go files in preparation for a commit

# Get the go files to format
GO_FILES=$(git diff --cached --name-only | grep ".go$")
TOPLEVEL=$(git rev-parse --show-toplevel)

if [[ "$GO_FILES" = "" ]]; then
  exit 0
fi

GOFMT=$(which gofmt)

if [[ ! -x "$GOFMT" ]]; then
  printf "No gofmt available..."
  exit 1
fi

for FILE in $GO_FILES; do
  $GOFMT -w "$TOPLEVEL/$FILE"
  git add "$TOPLEVEL/$FILE"
done

exit 0
