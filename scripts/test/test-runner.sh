#!/bin/bash

source ./scripts/test/lib/libconst.sh

echo "DIR_ROOT : $DIR_ROOT"
echo "DIR_BASE : $DIR_BASE"
echo "DIR_LIB  : $DIR_LIB"
echo

if [ ! -d "$DIR_BASE" ]
then
	echo "Execute this in the root-dir of goMS. The place where the 'scripts/' folder is."
	exit 1
fi

echo ">>> INIT"
CASES="reg_pos_once
reg_stress_multiple
send_pos_once
send_pos_multiple_messages
send_stress_multiple_clients
"

echo ">>> START"
echo

for CASE_ID in $CASES
do
	echo ">>> SETUP.....: $CASE_ID"
	source "$DIR_LIB/libconst.sh"
	source "$DIR_LIB/setup.sh" > "$DIR_BASE/setup.log" 2>&1
	mv "$DIR_BASE/setup.log" "$DIR_RES/setup.log"

	echo -n ">>> RUN.......: $CASE_ID"
	"$DIR_BASE/cases/$CASE_ID/test.sh"	\
		"$DIR_LIB/libtest.sh"							\
		"$DIR_LIB/libconst.sh"						\
		"$DIR_ROOT"												\
		"$CASE_ID"												\
		> "$DIR_RES/test.log" 2>&1
	if [ $? -eq 0 ]
	then
		echo " [ PASS ]"
	else
		echo " [ FAIL ]"
	fi

	echo ">>> TEAR DOWN.: $CASE_ID"
	source "$DIR_LIB/tear_down.sh" > "$DIR_RES/tear_down.log" 2>&1

	echo ">>> FINISHED..: $CASE_ID"
	echo
done

echo ">>> DONE"
