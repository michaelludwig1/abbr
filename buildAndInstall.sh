#!/bin/bash

go build abbrbin.go

mkdir -p ~/.local
mkdir -p ~/.local/bin
mkdir -p ~/.config
mkdir -p ~/.config/abbr

rm ~/.local/bin/abbrbin 2> /dev/null
cp ./abbrbin ~/.local/bin/abbrbin
chmod +x ~/.local/bin/abbrbin

rm ~/.local/bin/abbr.sh 2> /dev/null
cp ./abbr.sh ~/.local/bin/abbr.sh
chmod +x ~/.local/bin/abbr.sh

touch ~/.config/abbr/aliases

if ! grep -q "alias abbr" ~/.bashrc; then
    cat bashrcentry >> ~/.bashrc
fi
. ~/.bashrc

