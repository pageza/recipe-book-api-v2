# --------------------------------------------------------
# Custom Devcontainer Shell Enhancements
# --------------------------------------------------------

# Ensure the terminal supports 256 colors.
export TERM=xterm-256color

# Enable color support for ls and add ls aliases.
alias ls='ls --color=auto'
alias ll='ls -alF'
alias la='ls -A'
alias l='ls -CF'

# Enable colorized grep output.
alias grep='grep --color=auto'
alias fgrep='fgrep --color=auto'
alias egrep='egrep --color=auto'

# Set a nice, informative prompt.
export PS1="\[\033[1;32m\]\u@\h:\[\033[0m\]\w \$ "

# Enable tab autocomplete from bash-completion if available.
if [ -f /usr/share/bash-completion/bash_completion ]; then
  . /usr/share/bash-completion/bash_completion
elif [ -f /etc/bash_completion ]; then
  . /etc/bash_completion
fi

# Docker compose shortcuts.
alias dcbn="docker compose build --no-cache"                           # Build without cache.
alias dcup="docker compose up"                                           # Start containers.
alias dcbnup="docker compose build --no-cache && docker compose up"      # Build (no cache) then start.

# Go application shortcuts.
alias gotest="go test -v ./..."            # Run all tests with verbose output.
alias gobuild="go build ./..."             # Build all Go packages.
alias gofmt="gofmt -w ."                   # Format all Go files (recursively).

# Other convenience settings can be added below.
# For example, you might want to set your favorite editor:
export EDITOR=vim 