#!/usr/bin/env bash
set -euo pipefail

# Run one initialization step after the SQL purge, with the API on Sepolia.
# Async responses do not prove completion. Check API logs before the next step.
# Never run against an existing platform: these endpoints are not all idempotent.

if [[ $# -ne 1 ]]; then
    printf 'Usage: bash tools/db/deploy_fresh_platform.sh <1|2|3|4|5>\n' >&2
    exit 2
fi

case "$1" in
    1) path='/trex/init/implementations'; expected='202' ;;
    2) path='/trex/init/contract/factory/identity'; expected='200' ;;
    3) path='/trex/init/authority/configure'; expected='200' ;;
    4) path='/trex/init/contract/factory/trexFactory'; expected='202' ;;
    5) path='/trex/init/suite/deploy'; expected='202' ;;
    *) printf 'Step must be 1, 2, 3, 4 or 5.\n' >&2; exit 2 ;;
esac

: "${API_BASE_URL:?Set API_BASE_URL including /api/v1}"
: "${ACCESS_TOKEN:?Set ACCESS_TOKEN to a new ADMIN or MANAGER access token}"

case "$API_BASE_URL" in
    https://*/api/v1|https://*/api/v1/) ;;
    *) printf 'API_BASE_URL must use HTTPS and end with /api/v1.\n' >&2; exit 2 ;;
esac

# Feed authorization through stdin so the token is not put in curl's argv.
case "$ACCESS_TOKEN" in
    *[!a-zA-Z0-9._-]*)
        printf 'ACCESS_TOKEN must be a JWT without whitespace.\n' >&2
        exit 2
        ;;
esac

response=$(printf 'header = "Authorization: Bearer %s"\n' "$ACCESS_TOKEN" |
    curl --config - --silent --show-error --connect-timeout 10 --max-time 600 \
        --request POST --header 'Accept: application/json' \
        --write-out $'\n%{http_code}' "${API_BASE_URL%/}${path}")
status=${response##*$'\n'}
body=${response%$'\n'*}
printf '%s\nHTTP %s\n' "$body" "$status"

if [[ "$status" != "$expected" ]]; then
    printf 'Unexpected response; stop here and inspect the API logs. Do not retry blindly.\n' >&2
    exit 1
fi
if [[ "$expected" == '202' ]]; then
    printf 'Accepted, NOT completed. Wait for job success in the API logs before continuing.\n'
else
    printf 'Step completed. Check the response before continuing.\n'
fi
