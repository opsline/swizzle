#!/bin/bash

# Run chalk to load secrets early
# (it's ok if this fails)
eval $(chalk $CHALK_OPTS)

/usr/local/bin/swizzle
