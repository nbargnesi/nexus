#/usr/bin/env bash

# The next three lines are for the go shell.
export SCRIPT_HELP="Run go build everytime files are modified."
export SCRIPT_DESC="build-loop"
[[ "${BASH_SOURCE[0]}" != "${0}" ]] && return 0

# Normal script execution starts here.
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"/../
source "$DIR"/env.sh || exit 1
cd "$DIR" || exit 1

CMD="go build -x"
which colorgo >/dev/null
if [ $? -eq 0 ]; then
    CMD="colorgo build -x"
fi
while true; do
    gate || exit 1
    clear
    echo "[BUILDING]"
    $CMD && echo -e "\n[COMPLETE]"
done
exit 0

