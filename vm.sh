#!/bin/bash

set -eu

if [ $# -eq 0 ]; then
  ./vm.sh help
  exit 1
fi

readonly CMD=$1
readonly MUTAGEN_FILE=mutagen.yml
readonly MUTAGEN_LOCK_FILE=mutagen.yml.lock
readonly VM_NAME=ps.vagrant

start() {
  echo "Starting up virtual machine..."

  if vagrant status ${VM_NAME} | grep "running (virtualbox)" > /dev/null 2>&1; then
    echo "Virtual machine is already running."
    exit 1
  fi
  vagrant up

  test -e "${MUTAGEN_LOCK_FILE}" && mutagen project terminate
  mutagen project start -f "${MUTAGEN_FILE}"

  echo "Virtual machine has started. ✨"
}

stop() {
  echo "Shutting down virtual machine..."

  test -e "${MUTAGEN_LOCK_FILE}" && mutagen project terminate

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

  mutagen project reset -f "${MUTAGEN_FILE}"

  echo "Virtual machine has restarted. 👍"
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
    echo "Command not supported."
    exit 1
  }
esac
