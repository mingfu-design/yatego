#!/bin/bash

echo "Starting Yate"

/opt/yate/startyate.sh

echo "Starting demo http srv"

cd /vagrant/cmd/http-callback
./http-callback&
