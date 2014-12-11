#!/usr/bin/env bash
DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Pull in standard functions, e.g., default.
source "$DIR/.gosh.sh" || return 1
default CUSTOM_ENV_SH "$DIR/env.sh.custom"
assert_source "$CUSTOM_ENV_SH" || return 1

### GENERAL ENV VARS ###
default DIR             "$DIR"
default GL_VERSION      "master"

### GREENLINE DEFAULT PORTS ###
default GL_BCAST_INGRESS_PORT   9002
default GL_BCAST_EGRESS_PORT    9003
default GL_RR1_INGRESS_PORT     9004
default GL_RR1_EGRESS_PORT      9005
default GL_RR2_INGRESS_PORT     9006
default GL_RR2_EGRESS_PORT      9007

### PATHS ###
default BUILD           "$DIR"/build

### GOLANG ###
default GL_BUILD_ARGS   "-o $BUILD/greenline"
default GL_INSTALL_ARGS ""

### THE GO SHELL ###
default GOSH_SCRIPTS    "$DIR"/scripts
default GOSH_PROMPT     "gosh \e[0;32mgreenline\e[0m (?|#|#?)> "

