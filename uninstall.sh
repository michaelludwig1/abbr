#!/bin/bash

rm ~/.local/bin/abbrbin 2> /dev/null
rm ~/.local/bin/abbr.sh 2> /dev/null
rm ~/.config/abbr/aliases 2> /dev/null
rm ~/.config/abbr/aliastmp 2> /dev/null

if [ -d ~/.config/abbr ]; then
    if [ ! "$(ls -A ~/.config/abbr)" ]; then
         rmdir ~/.config/abbr
    fi
fi

startLine=$(grep -xn "# The following is for the abbr program" ~/.bashrc | cut -d : -f 1)
if ! test -z "$startLine"; then
    endLine=$((startLine+4))
    sed -i "${startLine},${endLine}d" ~/.bashrc
fi

. ~/.bashrc

