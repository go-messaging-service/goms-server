#!/bin/bash

DIR_ROOT="$3"
CASE_ID="$4"
source "$2"

function assert_no_errors()
{
	amount_server=$(grep -i "error" "$DIR_RES/server.log" | wc -l)
	if [ $amount_server -ne 0 ]
	then
		fail
	fi

	amount_client=$(grep -i "error" "$DIR_RES/test.log" | wc -l)
	if [ $amount_client -ne 0 ]
	then
		fail
	fi
}

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

	amount_client=$(grep -E "^\{\"messagetype\":\"message\"" "$DIR_RES/test.log" | grep -E "\"data\":\"$1\"\}$" | wc -l)
	if [ $amount_client -ne 1 ]
	then
		fail
	fi
}

# $1 - The content of the message
# $2 - The amount of times sent
function assert_sent_n()
{
	amount_server=$(grep -E "send message with data: " "$DIR_RES/server.log" | wc -l)
	if [ $amount_server -ne $2 ]
	then
		fail
	fi
}

# $1 - The content of the message
# $2 - The amount of times received
function assert_received_n()
{
	amount_client=$(grep -E "^\{\"messagetype\":\"message\"" "$DIR_RES/test.log" | grep -E "\"data\":\"$1\"\}$" | wc -l)
	echo $amount_client
	if [ $amount_client -ne $2 ]
	then
		fail
	fi
}

function fail()
{
	echo "FAIL"
	echo "Stack trace:"
	echo "    0 : ${BASH_SOURCE[1]} : ${BASH_LINENO[0]}  >  $(head -n ${BASH_LINENO[0]} ${BASH_SOURCE[1]} | tail -n 1)"
	echo "    1 : ${BASH_SOURCE[2]} : ${BASH_LINENO[1]}  >  $(head -n ${BASH_LINENO[1]} ${BASH_SOURCE[2]} | tail -n 1)"
	exit 1
}

function wait_tiny()
{
	sleep 0.1
}

function wait_short()
{
	sleep 0.5
}

function wait_med()
{
	sleep 1
}

function wait_long()
{
	sleep 2
}
