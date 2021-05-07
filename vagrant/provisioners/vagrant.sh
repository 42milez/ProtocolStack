#!/bin/bash

readonly BASH_PROFILE=~/.bash_profile
readonly BASHRC=~/.bashrc

#  Homebrew
# --------------------------------------------------
if ! type brew > /dev/null 2>&1; then
  CI=true /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
  test -d /home/linuxbrew/.linuxbrew && eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
  test -r ~/.bash_profile && echo "eval \$($(brew --prefix)/bin/brew shellenv)" >> "${BASH_PROFILE}"
  echo "" >> "${BASH_PROFILE}"
fi

#  Go
# --------------------------------------------------
if ! type go > /dev/null 2>&1; then
  brew install go
  # shellcheck disable=SC2016
  {
    echo ""
    echo "# Go"
    echo "export GO111MODULE=on"
    echo 'export GOBIN="${HOME}/.bin"'
    echo 'export GOMODCACHE="${HOME}/.cache/go_mod"'
    echo 'export GOPATH="${HOME}/.go"'
    echo 'PATH="${PATH}:${GOBIN}"'
    echo "export PATH"
    echo ""
  } >> "${BASHRC}"
  # shellcheck disable=SC1090
  . "${BASHRC}"
fi

#  Go - Delve
# --------------------------------------------------
if ! type dlv > /dev/null 2>&1; then
  go install github.com/go-delve/delve/cmd/dlv@latest
fi

#  Go - golangci-lint
# --------------------------------------------------
if ! type golangci-lint > /dev/null 2>&1; then
  brew install golangci-lint
fi
