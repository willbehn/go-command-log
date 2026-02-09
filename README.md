# go-comman-log

A lightweight command-line tool that records shell commands to a local SQLite database. This is a work in progress, but the goal is to offer som stats on your shell habits, and to make it easy to search for and use previous used commands. 

The CLI is created using golang with cobra, which is a first, so feel free to fork and add pull requests/ suggestions.

> Heads up: After enabling the hook in `.zshrc`-file all terminal commands will be logged. The SQLite database is currently **not encrypted** in any way, so be careful not to log anything secret or sensitive. If you start a command with **space** the command will not be logged. Specific keywords and phrases can be added to an ignore list, see `scripts/crumb-preexec.zsh`.


## Demo

![Made with VHS](https://vhs.charm.sh/vhs-5uHwmJTSilN6JYZbsxwpGQ.gif)

## Prerequisites
- Go installed
- Zsh

## Install 
> The install is currently very convoluted and wierd, will fix soon

Clone the repository
```bash
git clone https://github.com/yourname/go-command-log
```

Build the binary 
```bash
mkdir -p "<absolute-file-path>/go-command-log/bin" "<absolute-file-path>/go-command-log/database"
cd "<absolute-file-path>/go-command-log"
go build -o bin/commandlog
```

## Enable shell logging
Add this to **~/.zshrc**:
```bash
# Crumb
export CL_BIN="<absolute-file-path>go-command-log/bin/commandlog"
export CL_DB="<absolute-file-path>/go-command-log/database/history.db"

export PATH="<absolute-file-path>/go-command-log/bin:$PATH"

source "<absolute-file-path>/go-command-log/scripts/crumb-preexec.zsh"

# (Optional) short alias
alias cl=$CL_BIN
```

Then reload your shell:
```bash
source ~/.zshrc
```

Run `crumb init` to initalise database with schema

```bash
commandlog init
```

## Security and privacy

- Database is stored at `$CL_DB` and is **unencrypted**.
To avoid logging a command: start the line with a space, or add patterns to the ignore list in `scripts/crumb-preexec.zsh`.
Consider adding `database/history.db` to your .gitignore.

## TODO
- [ ] Make search fuzzy
- [ ] Add `crumb stats`commnd + flags
- [ ] Add more support for flags
- [ ] Autocompletion helpers, `crumb run X`
- [ ] Make install easier
- [ ] Make it work with different shells
- [ ] Move where database is stored
