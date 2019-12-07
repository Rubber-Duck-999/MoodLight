#!/bin/sh

echo "Starting Node"

export BWD=/home/simon/Documents/HouseGuard

if [ -e $BWD/../logs/faulthandler.txt ];
then
  rm -f $BWD/../logs/faulthandler.txt
fi

if [ ! -d $BWD/../logs ];
then
    mkdir $BWD/../logs
fi
$BWD/FaultHandler/bin/exeFaultHandler -f $BWD/FaultHandler/config.yml > $BWD/../logs/faulthandler.txt 2>&1 &
