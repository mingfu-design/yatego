#!/bin/bash

BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

glide install

echo "Building mock api..."
cd $BASEDIR/tools/api-mock
go build

cd $BASEDIR/cmd/
echo "Building IVRs ..."
for d in */ ; do
    cd $BASEDIR/cmd/$d
    go build
done

echo "Starting dev. env. ..."

cd $BASEDIR
./start_dev.sh