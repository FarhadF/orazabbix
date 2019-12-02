#!/bin/sh
while true
do
    while IFS="" read -r p || [ -n "$p" ]
    do
        ./orazabbix $p
    done < bases.txt
	sleep 30
done
