#!/bin/bash

go build abbrbin.go
rm ~/.local/bin/abbrbin 2> /dev/null
cp ./abbrbin ~/.local/bin/abbrbin
chmod +x ~/.local/bin/abbrbin

rm ~/.local/bin/abbr 2> /dev/null
cp ./abbr ~/.local/bin/abbr
chmod +x ~/.local/bin/abbr

mkdir -p ~/.config/abbr
touch ~/.config/abbr/aliases

if ! grep -q "alias abbr" ~/.bashrc; then
    cat bashrcentry >> ~/.bashrc
fi
. ~/.bashrc

