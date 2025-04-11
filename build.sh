#!/bin/bash

# Define the output binary name
OUTPUT_NAME=disc-go-bot

# Define the output directory
DIST_DIR="dist"

# Create dist directory if it doesn't exist
mkdir -p $DIST_DIR

# Define the platforms and architectures you want to build for
PLATFORMS="linux/amd64 windows/amd64 darwin/amd64"

# Check if we should build Windows as GUI application
BUILD_GUI=0
if [ "$1" == "--gui" ]; then
  BUILD_GUI=1
fi

# Check if building only for current platform
BUILD_CURRENT_ONLY=0
if [ "$1" == "--current" ] || [ "$2" == "--current" ]; then
  BUILD_CURRENT_ONLY=1
fi

# Get current platform
CURRENT_OS=$(go env GOOS)
CURRENT_ARCH=$(go env GOARCH)
echo "Current platform: $CURRENT_OS/$CURRENT_ARCH"

if [ $BUILD_CURRENT_ONLY -eq 1 ]; then
  echo "Building only for current platform: $CURRENT_OS/$CURRENT_ARCH"
  PLATFORMS="$CURRENT_OS/$CURRENT_ARCH"
fi

# Iterate over the platforms and build the binary for each
for platform in $PLATFORMS; do
  GOOS=${platform%/*}
  GOARCH=${platform#*/}
  output_name="${DIST_DIR}/${OUTPUT_NAME}_${GOOS}_${GOARCH}"

  build_flags=""

  echo "Building for $GOOS/$GOARCH..."

  # Add .exe extension for Windows builds
  if [ "$GOOS" == "windows" ]; then
    output_name+=".exe"

    # Add GUI flag for Windows if requested
    if [ $BUILD_GUI -eq 1 ]; then
      echo "Building with GUI flag..."
      export CGO_ENABLED=1
      build_flags="-ldflags=-H=windowsgui"
    else
      export CGO_ENABLED=0
    fi
  elif [ "$GOOS" == "$CURRENT_OS" ]; then
    # Building for current OS, CGO should work
    export CGO_ENABLED=1
  else
    # Cross-compilation, disable CGO
    export CGO_ENABLED=0
  fi

  # Build the binary
  env GOOS=$GOOS GOARCH=$GOARCH go build $build_flags -o $output_name

  if [ $? -eq 0 ]; then
    echo "Built binary: $output_name"
  else
    echo "Error building for $GOOS/$GOARCH"
  fi
done

# Also create a copy of the current platform binary in the dist root with the basic name
if [ $BUILD_CURRENT_ONLY -eq 0 ]; then
  echo "Creating shortcut binary for current platform..."
  if [ "$CURRENT_OS" == "windows" ]; then
    cp "${DIST_DIR}/${OUTPUT_NAME}_${CURRENT_OS}_${CURRENT_ARCH}.exe" "${DIST_DIR}/${OUTPUT_NAME}.exe"
  else
    cp "${DIST_DIR}/${OUTPUT_NAME}_${CURRENT_OS}_${CURRENT_ARCH}" "${DIST_DIR}/${OUTPUT_NAME}"
  fi
fi