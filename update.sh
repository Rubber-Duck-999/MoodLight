#!/bin/sh

cd $HOME/Documents/HouseGuard-FaultHandler/src

git pull

go clean

go build

if [ -f exeFaultHandler ];
then
    echo "FH File found"
    if [ -f $HOME/Documents/Temp/exeFaultHandler ];
    then
        echo "FH old removed"
        rm -f $HOME/Documents/Temp/exeFaultHandler
    fi
    mv exeFaultHandler $HOME/Documents/Temp/exeFaultHandler
fi