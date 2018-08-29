#!/bin/bash

# Params
# $1 - Case ID
CASE_ID=$1

source ./scripts/test/constants.sh

echo $CASE_ID

REG=$(cat "$DIR_CASE/messages/register.a-b" | tr -d "\n\t[:space:]")
echo "Send '$REG'"

# Send registration to server and wait a bit so that the server can handle it and print some logs
echo $REG | "$DIR_ROOT/connect.sh" &
sleep 0.5

# TODO Check if string was there. Create another .sh file with some assertions
grep -E ": Register a$" "$DIR_RES/server.log"
grep -E ": Register b$" "$DIR_RES/server.log"
