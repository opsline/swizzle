#!/bin/bash

# Run chalk to load secrets early
# (it's ok if this fails)
_ct=$(mktemp) ; rm -f $_ct ; chalk --debug > $_ct && source $_ct ; rm -f $_ct ; unset _ct

if [[ -n "$1" ]]; then
  set -x
  exec $@
fi

/usr/local/bin/swizzle

