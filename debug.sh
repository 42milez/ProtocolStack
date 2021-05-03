#!/bin/bash

set -eu

if [ $# -eq 0 ]; then
  ./debug.sh help
  exit 1
fi

readonly CMD=$1
readonly DELVE="dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec"
readonly TCP_CLIENT=./bin/tcp_client
readonly TCP_SERVER=./bin/tcp_server

case $CMD in
  "client") $DELVE $TCP_CLIENT;;
  "server") $DELVE $TCP_SERVER;;
  "help") {
    echo "Usage:"
    echo ""
    echo "  ./debug.sh <command>"
    echo ""
    echo "The commands are:"
    echo ""
    echo "  client   run the TCP client, and begin a debug session"
    echo "  server   run the TCP server, and begin a debug session"
  };;
  *) {
    echo "Command not supported."
    exit 1
  }
esac
