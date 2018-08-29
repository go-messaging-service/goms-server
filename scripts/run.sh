#!/bin/bash

echo "========="
echo "  CLEAN"
echo "========="
echo "Removing old build ..."
if [ -d ./bin/ ] && ! [ -z "$(ls -A ./bin/)" ]; then
	rm -r ./bin/*
	if [ $? -ne 0 ]; then
		echo ""
		echo "CLEANUP FAILED!"
		exit 1
	fi
fi
echo "CLEANUP DONE"
echo ""

echo "========="
echo " COMPILE"
echo "========="
echo "Calling './scripts/compile.sh"
sh ./scripts/compile.sh
if [ $? -ne 0 ]; then
	echo ""
	echo "COMPILING FAILED!"
	exit 2
fi
echo "COMPILING DONE"
echo ""

echo "========="
echo "  START"
echo "========="
./bin/goMS-server $@

# run "nc localhost 55545" or "sh connect.sh" to connect to the server via terminal
