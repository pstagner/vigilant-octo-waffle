#!/usr/bin/env bash
# source ./src/sourceror.bash

set -eu
if [[ -f ./.env ]]; then
  set -a && source ./.env && set -a
else
  echo 'no env'
  exit 1
fi

# Define a function to get the hosts file path
get_hosts_file_path() {
  if command -v wslpath &> /dev/null; then
    echo "/mnt/c/Windows/System32/drivers/etc/hosts"
  else
    echo "/etc/hosts"
  fi
}

# Get the hosts file path
HOSTS_FILE=$(get_hosts_file_path)

TMP=$(mktemp -d)
trap 'rm -Rf $TMP' EXIT
envsubst < "src/hosts" > "$TMP/hosts"

touched=0

while read -r line
do
  linerrr "$line" "$HOSTS_FILE"
  if [[ $? -eq 0 ]]; then
    ((touched++))
  fi
done < "$TMP/hosts"

if [[ ${touched} -eq 1 ]]; then
  echo "$HOSTS_FILE file has been updated"
else
  echo 'no changes were necessary'
fi