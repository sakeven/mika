#!/bin/bash

work_dir=$(dirname $0)
cd ${work_dir}

if [ ! -d "bin" ]; then
    mkdir bin
fi

for key in client server;
do
    echo building ${key} ...
    go build -o bin/${key} ./cmd/${key}
    if [ $? -ne 0 ]; then
        echo build ${key} failed!
        exit 1
    fi
done

echo success