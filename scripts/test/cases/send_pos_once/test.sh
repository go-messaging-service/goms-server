#!/bin/bash

source $1

REG=$(cat "$DIR_CASE/messages/register" | tr -d "\n")
SEND=$(cat "$DIR_CASE/messages/send" | tr -d "\n")
echo "Send '$REG'"
echo "Send '$SEND'"

# Register someone and wait a bit
echo "$REG" | "$DIR_ROOT/connect.sh" &
wait_tiny

# Send some data over to the client above
(
	echo "$REG"
	echo "$SEND"
) | "$DIR_ROOT/connect.sh" &
wait_tiny

assert_no_errors
assert_registered_n "a" 2
assert_registered_n "with spaces" 2
assert_sent_once "$MSG"
