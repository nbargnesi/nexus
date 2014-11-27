#!/bin/bash

# The next three lines are for the go shell.
export SCRIPT_NAME="release"
export SCRIPT_HELP="Create a release."
[[ "$GOGO_GOSH_SOURCE" -eq 1 ]] && return 0

# Normal script execution starts here.
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$DIR"/env.sh || exit 1
cd "$DIR" || exit 1

assert_env GH_KEY || exit 1
GH_API="https://api.github.com"

prompt_env GL_VERSION "Version (e.g., 0.1): "
git tag v$GL_VERSION || exit 1

RELEASE="
{
    \"tag_name\": \"v$GL_VERSION\",
    \"target_commitish\": \"master\",
    \"name\": \"v$GL_VERSION\",
    \"body\": \"Greenline binaries for Linux and OS X\",
    \"draft\": false,
    \"prerelease\": false
}
"

RESULT=$(echo $RELEASE | curl -su $GH_KEY:x-oauth-basic \
         $GH_API/repos/formwork-io/greenline/releases -d@-)
[[ $? != 0 ]] && exit 1
ID=$(echo $RESULT | jq '.id')
[[ $? != 0 ]] && exit 1
echo "Release id is: $ID"

