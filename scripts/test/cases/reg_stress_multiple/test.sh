#!/bin/bash

source $1

REG=$(cat "$DIR_CASE/messages/register.a-b" | tr -d "\n\t[:space:]")
echo "Send '$REG'"

# Send registration to server and wait a bit so that the server can handle it and print some logs
(
for i in {1..10000}
do
	echo $REG
done
) | "$DIR_ROOT/connect.sh" &
sleep 2

assert_registered_once "a"
assert_registered_once "b"
