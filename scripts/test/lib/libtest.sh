#!/bin/bash

DIR_ROOT="$2"
DIR_BASE="$3"
DIR_RES="$4"
DIR_CASE="$5"

# TODO Create an assert_no_err which will be needed very often

# $1 - The name of the topic
function assert_registered_once()
{
	amount=$(grep -E " Register $1\$" "$DIR_RES/server.log" | wc -l)
	if [ $amount -ne 1 ]
	then
		fail
	fi
}

# $1 - The name of the topic
# $2 - The amount of times registered
function assert_registered_n()
{
	amount=$(grep -E " Register $1\$" "$DIR_RES/server.log" | wc -l)
	if [ $amount -ne $2 ]
	then
		fail
	fi
}

# $1 - The content of the message
function assert_sent_once()
{
	amount_server=$(grep -E "send message with data: " "$DIR_RES/server.log" | wc -l)
	if [ $amount_server -ne 1 ]
	then
		fail
	fi
	
	amount_client=$(grep -E "\"messagetype\":\"message\"" "$DIR_RES/test.log" | grep -E "\"data\":\"$1\"" | wc -l)
	if [ $amount_server -ne 1 ]
	then
		fail
	fi
}

function fail()
{
	echo "FAIL: ${BASH_SOURCE[2]} : ${BASH_LINENO[1]}  >  $(head -n ${BASH_LINENO[1]} ${BASH_SOURCE[2]} | tail -n 1)"
	exit 1
}
