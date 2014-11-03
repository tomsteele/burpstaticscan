#!/bin/sh
 
for OS in "linux" "darwin" "windows"; do
    GOOS=$OS CGO_ENABLED=0 GOARCH=amd64 go build
    FOLDER=burpstaticscan-1.0.0$OS-amd64
    ARCHIVE=$FOLDER.tar.gz
    mkdir $FOLDER
    cp LICENSE $FOLDER
    cp README.md $FOLDER
    if [ $OS = "windows" ] ; then
        cp burpstaticscan.exe $FOLDER
        rm burpstaticscan.exe
    else
        cp burpstaticscan $FOLDER
        rm burpstaticscan
    fi
    tar -czf $ARCHIVE $FOLDER
    rm -rf $FOLDER
    echo $ARCHIVE
done
