#!/bin/bash

set -eu

readonly CMD=$1
readonly DELVE="dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec"
readonly TCP_CLIENT=./bin/tcp_client
readonly TCP_SERVER=./bin/tcp_server

case $CMD in
  "client") $DELVE $TCP_CLIENT;;
  "server") $DELVE $TCP_SERVER;;
  *) {
    echo "Command not supported."
    exit 1
  }
esac
