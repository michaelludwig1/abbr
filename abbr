#!/bin/bash

if test ! -t 0; then
    STDIN="$(cat /dev/stdin)"
else
    STDIN="$(history | tail -2 | head -1 | cut -c8-999)"
fi

echo $STDIN | ~/.local/bin/abbrbin "$@"

if [ -f ~/.config/abbr/aliastmp ]; then
    . ~/.config/abbr/aliastmp
    rm ~/.config/abbr/aliastmp
fi

