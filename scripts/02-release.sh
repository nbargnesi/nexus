#!/usr/bin/env bash
mkdir -p dist || exit 1

TARGETS="linux/amd64"
for target in "$TARGETS"; do
    GOOS=${target%/*}
    GOARCH=${target#*/}
    echo "Creating $GOOS $GOARCH binary..."
    GOOS=$GOOS GOARCH=$GOARCH go build -o dist/greenline
done

