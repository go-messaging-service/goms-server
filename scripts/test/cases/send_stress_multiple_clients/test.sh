#!/bin/bash

source $1

REG=$(cat "$DIR_CASE/messages/register" | tr -d "\n")
SEND=$(cat "$DIR_CASE/messages/send" | tr -d "\n")
echo "Send '$REG'"
echo "Send '$SEND'"

for i in {1..100}
do
	# Register someone and wait a bit
	echo "$REG" | "$DIR_ROOT/connect.sh" &
done
wait_med

# Send some data over to the client above
(
	echo "$REG"
	echo "$SEND"
) | "$DIR_ROOT/connect.sh" &
wait_tiny

assert_no_errors
assert_registered_n "a" 101
assert_registered_n "with spaces" 101
assert_sent_n "$MSG" 1
assert_received_n "$MSG" 100
