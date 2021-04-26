#!/bin/bash

set -eu

readonly CMD=$1
readonly VM_NAME=ps.vagrant

start() {
  if ! vagrant status ${VM_NAME} | grep "running (virtualbox)" > /dev/null 2>&1; then
    echo "already started."
    exit 1
  fi
  vagrant up
  mutagen project start -f mutagen.yml
}

stop() {
  if ! vagrant status ${VM_NAME} | grep "poweroff (virtualbox)" > /dev/null 2>&1; then
    echo "already stopped."
    exit 1
  fi
  mutagen project terminate
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
