#/bin/bash
set -e
# This script retrieves the JWT keys from the runtime directory and outputs them in JSON format.
# It is used by Terraform to allow it to generate the keys and template them in the auth Helm chart.
privateKey=$(cat .runtime/privateKey.json | base64 -w 0)
publicKey=$(cat .runtime/publicKey.json | base64 -w 0)
echo "{\"privateKey\": \"$privateKey\", \"publicKey\": \"$publicKey\"}"
# The output is a compact JSON object containing the base64 encoded private and public keys.
# This format is suitable for use in Terraform or other tools that require JSON input.
