#!/bin/bash

# Define the output binary name
OUTPUT_NAME=disc-go-bot

# Define the platforms and architectures you want to build for
PLATFORMS="linux/amd64 windows/amd64 darwin/amd64"

# Check if we should build Windows as GUI application
BUILD_GUI=0
if [ "$1" == "--gui" ]; then
  BUILD_GUI=1
fi

# Iterate over the platforms and build the binary for each
for platform in $PLATFORMS; do
  GOOS=${platform%/*}
  GOARCH=${platform#*/}
  output_name="${OUTPUT_NAME}_${GOOS}_${GOARCH}"

  build_flags=""

  # Add .exe extension for Windows builds
  if [ "$GOOS" == "windows" ]; then
    output_name+=".exe"
    # Add GUI flag for Windows if requested
    if [ $BUILD_GUI -eq 1 ]; then
      build_flags="-ldflags=\"-H=windowsgui\""
    fi
  fi

  echo "Building for $GOOS/$GOARCH..."
  if [ -n "$build_flags" ]; then
    eval env GOOS=$GOOS GOARCH=$GOARCH go build $build_flags -o $output_name
  else
    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name
  fi
  echo "Built binary: $output_name"
done