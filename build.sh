#!/bin/bash

# Define the output binary name
OUTPUT_NAME=disc-go-bot

# Define the platforms and architectures you want to build for
PLATFORMS="linux/amd64 windows/amd64 darwin/amd64"

# Iterate over the platforms and build the binary for each
for platform in $PLATFORMS; do
  GOOS=${platform%/*}
  GOARCH=${platform#*/}
  output_name="${OUTPUT_NAME}_${GOOS}_${GOARCH}"

  # Add .exe extension for Windows builds
  if [ "$GOOS" == "windows" ]; then
    output_name+=".exe"
  fi

  echo "Building for $GOOS/$GOARCH..."
  env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
  echo "Built binary: $output_name"
done