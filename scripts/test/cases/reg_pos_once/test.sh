#!/bin/bash

source $1

REG=$(cat "$DIR_CASE/messages/register.a-b" | tr -d "\n\t[:space:]")
echo "Send '$REG'"

# Send registration to server and wait a bit so that the server can handle it and print some logs
echo $REG | "$DIR_ROOT/connect.sh" &
sleep 0.5

assert_registered_once "a"
assert_registered_once "b"
