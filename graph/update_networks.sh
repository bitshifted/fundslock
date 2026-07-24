#!/bin/bash
NETWORK_NAME=$1
ADDRESS=$2
BLOCK=$3
FILE="networks.json"

# Ensure the current directory is graph/
# (Assuming the script is run from the blockchain root, this might need adjustment)
# Let's adjust to be relative to where the script is located (graph/)

# Ensure the directory exists
mkdir -p .

# Initialize if not exists
if [ ! -f "$FILE" ]; then echo '{}' > "$FILE"; fi

# Use jq to update the JSON
jq --arg net "$NETWORK_NAME" \
   --arg addr "$ADDRESS" \
   --arg block "$BLOCK" \
   '.[$net].FundsLock = {address: $addr, startBlock: $block | tonumber}' \
   "$FILE" > "$FILE.tmp" && mv "$FILE.tmp" "$FILE"
