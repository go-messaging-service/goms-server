#!/bin/bash

echo ">>> INIT"

source ./scripts/test/constants.sh
CASES="0 1"

echo ">>> START"
echo

for CASE_ID in $CASES
do
	echo ">>> SETUP $CASE_ID"
	"$DIR_BASE/setup.sh" "$CASE_ID"

	echo ">>> RUN $CASE_ID"
	"$DIR_BASE/cases/$CASE_ID/test.sh"
	
	echo ">>> FINISHED $CASE_ID"
	echo
done

echo ">>> DONE"
