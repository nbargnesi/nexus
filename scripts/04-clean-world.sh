#!/usr/bin/env bash

# The next three lines are for the go shell.
export SCRIPT_HELP="Cleans the tree of every artifact."
export SCRIPT_DESC="clean-world"
[[ "${BASH_SOURCE[0]}" != "${0}" ]] && return 0

# Normal script execution starts here.
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
cd "${DIR}" || exit 1
. "$DIR"/env.sh || exit 1

echo -en "Cleaning builds... "
rm -fr build
echo "ok"

