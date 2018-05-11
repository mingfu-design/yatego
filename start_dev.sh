#!/bin/bash

echo "Starting vagrant..."
vagrant up

echo "Starting Mock API..."

printf "\n\n\033[32;1m Test Call using softphone account 41587000201@172.28.128.3, pass: milan to number 925 \033[0m\n\n"

cd tools/api-mock
./api-mock