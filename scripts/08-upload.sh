#!/bin/bash

# The next three lines are for the go shell.
export SCRIPT_NAME="upload"
export SCRIPT_HELP="Upload a $(uname) $(uname -m) archive to a release."
[[ "$GOGO_GOSH_SOURCE" -eq 1 ]] && return 0

# Normal script execution starts here.
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$DIR"/env.sh || exit 1
cd "$DIR" || exit 1

assert_env GH_KEY || exit 1
assert_env DIST || exit 1
assert_env GL_RELEASE_ID || exit 1
prompt_env GL_VERSION "Version (e.g., 0.1): "
GH_UPLOAD_API="https://uploads.github.com"
GH_API_PATH="repos/formwork-io/greenline/releases/$GL_RELEASE_ID/assets"

ARTIFACT="greenline-v$GL_VERSION-$(uname)-$(uname -m)"
ARTIFACT_ROOT="$DIST"/$ARTIFACT
mkdir -p $ARTIFACT_ROOT
go build -o $ARTIFACT_ROOT/greenline || exit 1
cd $DIST || exit 1
tar caf $ARTIFACT.tar.gz $ARTIFACT || exit 1
HTTP_STATUS=$(curl --silent \
     -X POST \
     -u $GH_KEY:x-oauth-basic \
     -H "Content-Type: application/gzip" \
     -w %{http_code} \
     -o /dev/null \
     "$GH_UPLOAD_API/$GH_API_PATH?name=$ARTIFACT.tar.gz" \
     -F "filecomment=This is an image file" -F "name=@$ARTIFACT.tar.gz")
if [ "$HTTP_STATUS" != "201" ]; then
    echo "FAIL, status: $HTTP_STATUS"
    exit 1
fi

#     -d@$ARTIFACT.tar.gz)
