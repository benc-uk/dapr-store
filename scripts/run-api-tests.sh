#!/bin/bash

#
# Runs API tests using Postman collection and newman runner
#

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Check newman is installed
which newman > /dev/null || { echo -e "💥 Error! Command newman not installed, run 'npm install --global newman'"; exit 1; }

ENVIRONMENT_FILE=${1:-"$DIR/../testing/postman_env_local.json"}
TEST_FILE=${2:-"$DIR/../testing/postman_api_tests.json"}

newman run -e "$ENVIRONMENT_FILE" "$TEST_FILE"
