#!/bin/sh

echo "========="
echo " COMPILE"
echo "========="
bash ./scripts/compile.sh
echo ""
echo "========="
echo "  START"
echo "========="
./bin/goMS-server

# run "nc localhost 55545" or "sh connect.sh" to connect to the server via terminal
