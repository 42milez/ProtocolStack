#!/bin/bash

readonly CMD=$1
readonly DELVE="dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec"

case $CMD in
  "server") $DELVE ./bin/tcp_server;;
  "client") $DELVE ./bin/tcp_client;;
  *) {
    echo "Command not supported."
    exit 1
  }
esac
