#!/bin/bash

echo ">>> INIT"

source ./scripts/test/constants.sh
CASES="0"

echo ">>> START"
echo

for CASE_ID in $CASES
do
	echo ">>> SETUP $CASE_ID"
	"$DIR_BASE/setup.sh" "$CASE_ID"

	echo ">>> RUN $CASE_ID"
	"$DIR_BASE/cases/$CASE_ID/test.sh" "$CASE_ID"

	echo ">>> TEAR DOWN $CASE_ID"
	"$DIR_BASE/tear_down.sh"

	echo ">>> FINISHED $CASE_ID"
	echo
done

echo ">>> DONE"
