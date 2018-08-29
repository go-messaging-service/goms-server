#/bin/bash

# Params
# $1 - Case ID
CASE_ID=$1

source ./scripts/test/constants.sh

if [ ! -d "$DIR_BASE" ]
then
	echo "Execute this in the root-dir of goMS. The place where the 'scripts/' folder is."
	exit 1
fi

echo ">>> RESET $DIR_RES"
rm -rf "$DIR_RES"
mkdir "$DIR_RES"

echo ">>> LAUNCH SERVER"
nohup "$DIR_ROOT/run.sh" -c "$DIR_CASE/conf/server.json" 2>1 > "$DIR_CASE/server.log" &

# Wait a bit for the server to come up
for WAIT in $SRV_WAIT
do
	echo -ne "."
	sleep $WAIT
done
echo "[OK]"
