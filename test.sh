#!/bin/bash

readonly S="refs/heads/branch_name"
readonly LINK=$(
if [[ ${S} =~ ^"refs/heads/" ]]; then
  echo "exists"
else
  echo "not exists"
fi
)

echo "a"
echo "${LINK}"
