#!/bin/bash

echo ">>> INIT"
CASES="0"

echo ">>> START"
echo

for CASE_ID in $CASES
do
	echo ">>> SETUP $CASE_ID"
	source ./scripts/test/constants.sh
	source "$DIR_BASE/setup.sh"

	echo ">>> RUN $CASE_ID"
	source "$DIR_BASE/cases/$CASE_ID/test.sh"

	echo ">>> TEAR DOWN $CASE_ID"
	source "$DIR_BASE/tear_down.sh"

	echo ">>> FINISHED $CASE_ID"
	echo
done

echo ">>> DONE"
