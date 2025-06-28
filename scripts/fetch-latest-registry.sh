#!/bin/bash

UPSTREAM_REPO="https://github.com/LIFX/products.git"
CLONE_DIR="$(mktemp -d)"
SRC_DIR=$(git rev-parse --show-toplevel)/cmd/registry-gen/src
COMMIT_FILE="$SRC_DIR/products_commit.txt"

echo "Cloning LIFX/products repo"

# Clone the upstream repo into a temp dir
git clone --depth 1 "$UPSTREAM_REPO" "$CLONE_DIR" --quiet

# Get the latest commit hash from the clone (main branch assumed)
cd "$CLONE_DIR"
LATEST_COMMIT=$(git rev-parse HEAD)

# Compare to currently recorded commit
if [ -f "$COMMIT_FILE" ]; then
    CURRENT_COMMIT=$(cat "$COMMIT_FILE")
    if [ "$CURRENT_COMMIT" = "$LATEST_COMMIT" ]; then
        echo "Product registry commit is up to date ($CURRENT_COMMIT). No regeneration needed."
        rm -rf "$CLONE_DIR"
        exit 1  # ⚠️ exit non-zero to signal "no update needed"
    fi
fi

echo "Product registry updated: $LATEST_COMMIT"
# Copy new product registry and commit hash
cp products.json "$SRC_DIR/products.json"
echo "$LATEST_COMMIT" > "$SRC_DIR/products_commit.txt"

# Clean up
rm -rf "$CLONE_DIR"
