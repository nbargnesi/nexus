#!/usr/bin/env bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Pull in standard functions, e.g., default.
source "$DIR/.gosh.sh" || exit 1

### GENERAL ENV VARS ###
default DIR             "$DIR"
default CUSTOM_ENV      "$DIR/env.sh.custom"

### PATHS ###
default BUILD           "$DIR"/build

### GOLANG ###
default GL_BUILD_ARGS   "-o $BUILD/greenline"

### THE GO SHELL ###
default GOSH_SCRIPTS    "$DIR"/scripts
default GOSH_PROMPT     "gosh \e[0;32mgreenline\e[0m (?|#)> "

default CUSTOM_ENV  "$GL_DIR/env.sh.custom"

if [ -r "$CUSTOM_ENV_SH" ]; then
    source $CUSTOM_ENV_SH || exit 1
fi

