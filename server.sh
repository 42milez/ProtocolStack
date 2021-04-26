#!/bin/bash

set -eu

readonly CMD=$1
readonly VM_NAME=ps.vagrant

terminate_mutagen_project() {
  test -e mutagen.yml.lock && mutagen project terminate
}

start() {
  if vagrant status ${VM_NAME} | grep "running (virtualbox)" > /dev/null 2>&1; then
    echo "already started."
    exit 1
  fi
  vagrant up

  terminate_mutagen_project
  mutagen project start -f mutagen.yml
}

stop() {
  terminate_mutagen_project

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
