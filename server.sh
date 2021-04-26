#!/bin/bash

readonly CMD=$1

start() {
  vagrant up
  mutagen project start -f mutagen.yml
}

stop() {
  mutagen project terminate
  vagrant halt
}

if ! type mutagen > /dev/null 2>&1; then
  brew install mutagen-io/mutagen/mutagen
fi

case $CMD in
  "start") start;;
  "stop") stop;;
  *) {
    echo "unsupported command."
    exit 1
  }
esac
