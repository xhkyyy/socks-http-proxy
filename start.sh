#!/bin/bash

BASEDIR="$(
    cd "$(dirname "$0")"
    pwd -P
)"

"$BASEDIR/proxy" -f "$BASEDIR/cnf.json"