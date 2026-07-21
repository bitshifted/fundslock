#!/usr/bin/env bash
set -euo pipefail


if [ $# -lt 1 ]; then
    echo "Usage: $0 <kms-key-id>"
    exit 1
fi

export KMS_KEY_ID="$1"

# 2. Securely prompt for the MetaMask Key (Inputs are hidden from the screen)
echo -n "Enter MetaMask (or another wallet) Private Key (Hex format, without 0x): "
read -s METAMASK_HEX
echo "" # Prints a new line after hidden input

# Validate that something was entered
if [ -z "$METAMASK_HEX" ]; then
    echo "Error: Private key cannot be empty."
    exit 1
fi

echo "=== Starting Secure AWS KMS Key Import Workflow ==="

# 3. Setup a secure RAM disk or memory-only temporary folder
# Using /dev/shm (shared memory) on Linux ensures data never hits the hard drive
if [ -d "/dev/shm" ]; then
    TMP_DIR=$(mktemp -d /dev/shm/kms-import-XXXXXXXXXX)
else
    TMP_DIR=$(mktemp -d -t kms-import-XXXXXXXXXX)
fi
cd "$TMP_DIR"

# Strict cleanup handler to zero out files and unset memory variables
cleanup() {
    echo "=== Purging volatile memory and temp files ==="
    cd "$HOME"
    if command -v shred &> /dev/null; then
        shred -u -z -n 3 "$TMP_DIR"/* 2>/dev/null || true
    fi
    rm -rf "$TMP_DIR"
    # Zero out the sensitive shell variables
    METAMASK_HEX="0000000000000000000000000000000000000000000000000000000000000000"
    unset METAMASK_HEX
    echo "Security cleanup complete."
}
trap cleanup EXIT

# 4. Format the key strictly using memory pipes where possible
echo "Structuring and formatting key material..."
echo "302e0201010420${METAMASK_HEX}a00706052b8104000a" | xxd -r -p > private-key.raw

openssl ec -inform DER -in private-key.raw -out private-key.pem 2>/dev/null
openssl pkcs8 -topk8 -nocrypt -inform PEM -outform DER -in private-key.pem -out private-key.der

# 5. Fetch AWS KMS parameters
echo "Fetching AWS KMS wrapping parameters..."
aws kms get-parameters-for-import \
    --region "$AWS_REGION" \
    --key-id "$KMS_KEY_ID" \
    --wrapping-algorithm "RSAES_OAEP_SHA_256" \
    --wrapping-key-spec "RSA_2048" \
    --query '[PublicKey, ImportToken]' \
    --output json > import_params.json

jq -r '.[0]' import_params.json | base64 --decode > WrappingPublicKey.bin
jq -r '.[1]' import_params.json | base64 --decode > ImportToken.bin

# 6. Encrypt the key material 
echo "Encrypting key material..."
openssl pkeyutl -encrypt \
    -in private-key.der \
    -out EncryptedKeyMaterial.bin \
    -inkey WrappingPublicKey.bin \
    -keyform DER \
    -pubin \
    -pkeyopt rsa_padding_mode:oaep \
    -pkeyopt rsa_oaep_md:sha256 \
    -pkeyopt rsa_mgf1_md:sha256

# 7. Upload to AWS KMS
echo "Uploading to AWS KMS..."
aws kms import-key-material \
    --region "$AWS_REGION" \
    --key-id "$KMS_KEY_ID" \
    --encrypted-key-material fileb://EncryptedKeyMaterial.bin \
    --import-token fileb://ImportToken.bin \
    --expiration-model "KEY_MATERIAL_DOES_NOT_EXPIRE"

echo "=== Success! Key material imported securely. ==="
