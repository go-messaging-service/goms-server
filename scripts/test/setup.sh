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

echo "Reset $DIR_RES"
rm -rf "$DIR_RES"
mkdir "$DIR_RES"
