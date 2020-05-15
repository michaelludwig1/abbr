# abbr
Simple command line tool for managing abbreviations, i.e. aliases, for bash.

## Installation
1. Make sure, you have go installed
2. Download or clone the repo
3. Run buildAndInstall.sh

## Basic usage
The command is `abbr`. Run it with no parameters and you will see a help page. In summary the most important commands are:
* `abbr l` shows all aliases currently managed by abbr
* `abbr a name command` adds an alias `name` for `command`
* `abbr a name` can receive a string from stdin and will use `name` as an alias for it
* If stdin is empty, `abbr a name` adds an alias `name` for the last command executed, according to the history
* `abbr r name` removes the alias `name`


