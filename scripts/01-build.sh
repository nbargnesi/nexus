#/usr/bin/env bash

# The next three lines are for the go shell.
export SCRIPT_HELP="Run go build."
export SCRIPT_DESC="build"
[[ "${BASH_SOURCE[0]}" != "${0}" ]] && return 0

# Normal script execution starts here.
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$DIR"/env.sh || exit 1
cd "$DIR" || exit 1

CMD="go build"
which colorgo >/dev/null
if [ $? -eq 0 ]; then
    echo "[BUILDING]"
    CMD="colorgo build" && echo "[COMPLETE]"
fi
$CMD

