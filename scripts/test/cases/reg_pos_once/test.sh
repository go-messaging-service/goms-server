#!/bin/bash

source $1

REG=$(cat "$DIR_CASE/messages/register" | tr -d "\n")
echo "Send '$REG'"

# Send registration to server and wait a bit so that the server can handle it and print some logs
echo $REG | "$DIR_ROOT/connect.sh" &
wait_tiny

assert_no_errors
assert_registered_once "a"
assert_registered_once "with spaces"
