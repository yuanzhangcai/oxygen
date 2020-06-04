#!/bin/bash

count=`ps aux | grep 'oxygen' | grep -v 'restart.sh' | awk '{print $2}'`
echo $count

if [ -n "$count" ]; then
    ps aux | grep 'oxygen' | grep -v 'restart.sh' | awk '{print $2}' | xargs kill
    for i in $(seq 60)
    do
        count=`ps aux | grep 'oxygen' | grep -v 'restart.sh' | awk '{print $2}'`
        if [ -z "$count" ]; then
            break
        fi
        sleep 0.5
    done
fi

nohup /Users/zacyuan/MyWork/oxygen/oxygen > /dev/null 2>&1 &
