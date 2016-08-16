#!/bin/bash
set -euo pipefail

echo "== push"
env
if [ "$TYPE" == "note" ]; then
	echo 'I got a note!'
	cat
	echo
elif [ "$TYPE" == "url" ]; then
	echo "I should download $URL or something"
else
	echo "I don't know how to handle a $TYPE message"
fi
