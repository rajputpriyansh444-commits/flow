#!/bin/sh

# Set up constants and URLs
RELEASES_URL="https://github.com/rajputpriyansh444-commits/flow/releases"
FILE_BASENAME="flow"

# Parse CLI arguments
while [ "$#" -gt 0 ]; do
  case "$1" in
    --version)
      shift
      VERSION="$1"
      ;;
    --install-dir)
      shift
      INSTALL_DIR="$1"
      ;;
    *)
      ;;
  esac
  shift
done

# Check for curl and tar availability
command -v curl >/dev/null 2>&1 || { echo "curl is not installed. Aborting." >&2; exit 1; }
command -v tar >/dev/null 2>&1 || { echo "tar is not installed. Aborting." >&2; exit 1; }

echo "Downloading Flow version: $VERSION"

# Retrieve the latest version if not specified
if [ -z "$VERSION" ]; then
    VERSION="$(curl -sfL -o /dev/null -w %{url_effective} "$RELEASES_URL/latest" | rev | cut -f1 -d'/' | rev)"
fi

# Exit if VERSION is still empty
if [ -z "$VERSION" ]; then
    echo "Unable to get Flow version." >&2
    exit 1
fi

# Determine architecture
case "$(uname -m)" in
    x86_64)
        ARCHITECTURE="amd64"
        ;;
    aarch64)
        ARCHITECTURE="arm64"
        ;;
    *)
        ARCHITECTURE="$(uname -m)"
        ;;
esac

# Create a temporary directory for the download
TMPDIR=$(mktemp -d)
export TAR_FILE="${TMPDIR}/${FILE_BASENAME}_$(uname -s)_${ARCHITECTURE}.tar.gz"

# Download the binary
echo "Downloading Flow $VERSION..."
curl -sfLo "$TAR_FILE" "$RELEASES_URL/download/$VERSION/${FILE_BASENAME}_${VERSION:1}_$(uname -s)_${ARCHITECTURE}.tar.gz" || {
    echo "Failed to download the file. Check your internet connection and the URL." >&2
    exit 1
}

# Unpack and install
tar -xf "$TAR_FILE" -C "$TMPDIR" && sudo mv "${TMPDIR}/flow" /usr/local/bin/flow && sudo chmod +x /usr/local/bin/flow || {
    echo "Failed to install Flow." >&2
    exit 1
}

# Cleanup
rm -rf "$TMPDIR"
echo "Flow installed successfully and is available at /usr/local/bin/flow"

# Execute the binary with any provided arguments
"/usr/local/bin/flow" "$@"