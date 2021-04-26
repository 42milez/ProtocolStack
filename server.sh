#!/bin/bash

set -eu

readonly CMD=$1
readonly VM_NAME=ps.vagrant

start() {
  if vagrant status ${VM_NAME} | grep "running (virtualbox)" > /dev/null 2>&1; then
    echo "already started."
    exit 1
  fi
  vagrant up

  if [ -e mutagen.yml.lock ]; then
    mutagen project terminate
  fi
  mutagen project start -f mutagen.yml
}

stop() {
  if [ -e mutagen.yml.lock ]; then
    mutagen project terminate
  fi

  if vagrant status ${VM_NAME} | grep "poweroff (virtualbox)" > /dev/null 2>&1; then
    echo "already stopped."
    exit 1
  fi
  vagrant halt
}

case $CMD in
  "start") start;;
  "stop") stop;;
  *) {
    echo "unsupported command."
    exit 1
  }
esac
