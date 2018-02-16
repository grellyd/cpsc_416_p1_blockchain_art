#!/bin/bash

minerName=$1
localServer=$2

if [ -z $minerName ]; then
    echo "usage: ./start_miner.sh <MINER_NAME> [localServer]"
    echo ""
    exit 1
fi


if [ $localServer ]; then
  serverAddr="127.0.0.1:12345"
  echo "Calling ink-miner.go with server address $serverAddr"
else
  serverAddr="40.65.104.57:12345"
  echo "Calling ink-miner.go with server address $serverAddr"
fi

pubfile="keyfiles/$minerName.pub"
privfile="keyfiles/$minerName.priv"

go run ../ink-miner.go $serverAddr `cat $pubfile` `cat $privfile`

