#!/bin/bash

echo "::task~> Clean up & Build binaries files"
rm integration-test urlfounder 2>/dev/null
cd ../urlfounder
go build
mv urlfounder ../integration-test/urlfounder
cd ../integration-test
go build
echo "::done::"
echo "::task~> Run integration test"
./integration-test
echo "::done::"
if [ $? -eq 0 ]
then
  exit 0
else
  exit 1
fi
