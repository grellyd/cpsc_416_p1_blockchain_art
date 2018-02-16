#!/bin/bash

serverAddr="40.65.122.229:12345"
minerName=$1
pubfile="keyfiles/$minerName.pub"
privfile="keyfiles/$minerName.priv"

echo $pubfile
echo $privfile

go run ../ink-miner.go $serverAddr `cat $pubfile` `cat $privfile`

