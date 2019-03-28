#!/bin/bash

source $1

REG=$(cat "$DIR_CASE/messages/register" | tr -d "\n")
echo "Send '$REG'"

# Send registration to server and wait a bit so that the server can handle it and print some logs
(
for i in {1..10000}
do
	echo $REG
done
) | "$DIR_ROOT/connect.sh" &
wait_long

assert_registered_once "a"
assert_registered_once "with spaces"
