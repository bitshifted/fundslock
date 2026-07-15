#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: ./infra/kms-key-create.sh --policy-file <path> [options]

Create an empty external asymmetric AWS KMS key for later import of key material.

Options:
  --policy-file PATH   Path to a JSON policy document (required)
  --description TEXT   Description for the new key (default: Imported asymmetric KMS key)
  --alias NAME         Alias to create for the key (optional)
  --key-spec SPEC      KMS asymmetric key spec (default: ECC_SECP_256K1)
  --key-usage USAGE    KMS key usage (default: SIGN_VERIFY)
  -h, --help           Show this help message
EOF
}

policy_file=""
description="Imported asymmetric KMS key"
alias_name=""
key_spec="ECC_SECG_P256K1"
key_usage="SIGN_VERIFY"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --policy-file)
      policy_file="${2:-}"
      shift 2
      ;;
    --description)
      description="${2:-}"
      shift 2
      ;;
    --alias)
      alias_name="${2:-}"
      shift 2
      ;;
    --key-spec)
      key_spec="${2:-}"
      shift 2
      ;;
    --key-usage)
      key_usage="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [[ -z "$policy_file" ]]; then
  echo "Error: --policy-file is required." >&2
  usage >&2
  exit 1
fi

if [[ ! -f "$policy_file" ]]; then
  echo "Error: policy file not found: $policy_file" >&2
  exit 1
fi

if ! command -v aws >/dev/null 2>&1; then
  echo "Error: aws CLI is required and was not found in PATH." >&2
  exit 1
fi

policy_path="$(realpath "$policy_file")"

if [[ -z "${AWS_REGION:-}" ]]; then
  AWS_REGION="$(aws configure get region || true)"
fi

if [[ -z "${AWS_REGION:-}" ]]; then
  echo "Error: set AWS_REGION or configure your AWS CLI region before running this script." >&2
  exit 1
fi

key_id="$(aws kms create-key \
  --description "$description" \
  --origin EXTERNAL \
  --key-spec "$key_spec" \
  --key-usage "$key_usage" \
  --policy "file://$policy_path" \
  --query 'KeyMetadata.KeyId' \
  --output text \
  --region "$AWS_REGION")"

key_arn="$(aws kms describe-key \
  --key-id "$key_id" \
  --query 'KeyMetadata.Arn' \
  --output text \
  --region "$AWS_REGION")"

if [[ -n "$alias_name" ]]; then
  aws kms create-alias \
    --alias-name "alias/$alias_name" \
    --target-key-id "$key_id" \
    --region "$AWS_REGION" >/dev/null
fi

echo "Created KMS key: $key_id"
echo "Key ARN: $key_arn"
