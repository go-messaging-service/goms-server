#!/bin/bash

DIR_ROOT="$2"
DIR_BASE="$3"
DIR_RES="$4"
DIR_CASE="$5"

# $1 - The name of the topic
function assert_registered_once()
{
	amount=$(grep -E " Register $1\$" "$DIR_RES/server.log" | wc -l)

	if [ $amount -ne 1 ]
	then
		fail
	fi
}

function fail()
{
	echo "FAIL: ${BASH_SOURCE[2]} : ${BASH_LINENO[1]}  >  $(head -n ${BASH_LINENO[1]} ${BASH_SOURCE[2]} | tail -n 1)"
	exit 1
}
