#!/usr/bin/env bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Pull in standard functions, e.g., default.
source "$DIR/.gosh.sh" || return 1
default CUSTOM_ENV_SH "$DIR/env.sh.custom"
assert-source "$CUSTOM_ENV_SH" || return 1

### GENERAL ENV VARS ###
default DIR             "$DIR"
default CUSTOM_ENV_SH   "$DIR/env.sh.custom"

### THE GO SHELL ###
default GOSH_SCRIPTS    "$DIR"/scripts
default GOSH_CONTRIB    "$DIR"/scripts/gosh-contrib
default GOSH_PROMPT     "gosh \e[0;32mnexus request-protocol\e[0m (?|#|#?)> "

