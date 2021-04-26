#!/bin/bash

set -eu

readonly CMD=$1
readonly VM_NAME=ps.vagrant

start() {
  echo "Starting up virtual machine..."

  if vagrant status ${VM_NAME} | grep "running (virtualbox)" > /dev/null 2>&1; then
    echo "Virtual machine is already running."
    exit 1
  fi
  vagrant up

  test -e mutagen.yml.lock && mutagen project terminate
  mutagen project start -f mutagen.yml

  echo "Virtual machine has started!! ✨"
}

stop() {
  echo "Shutting down virtual machine..."

  test -e mutagen.yml.lock && mutagen project terminate

  if vagrant status ${VM_NAME} | grep "poweroff (virtualbox)" > /dev/null 2>&1; then
    echo "Virtual machine is already stopped."
    exit 1
  fi
  vagrant halt

  echo "Virtual machine has stopped. 😪"
}

restart() {
  echo "Restarting virtual machine..."

  if ! vagrant status ${VM_NAME} | grep "running (virtualbox)" > /dev/null 2>&1; then
    echo "Virtual machine is not running."
    exit 1
  fi
  vagrant reload

  mutagen project reset -f mutagen.yml

  echo "Virtual machine has restarted!! 👍"
}

case $CMD in
  "start") start;;
  "stop") stop;;
  "restart") restart;;
  "help") {
    echo "Usage:"
    echo ""
    echo "  ./vm.sh <command>"
    echo ""
    echo "The commands are:"
    echo ""
    echo "  start     start a VM and create a Mutagen session"
    echo "  stop      stop the VM and terminate the Mutagen session"
    echo "  restart   restart the VM and recreate a Mutagen session"
  };;
  *) {
    echo "unsupported command."
    exit 1
  }
esac
