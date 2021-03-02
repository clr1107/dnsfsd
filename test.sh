#!/bin/bash

echo "Running tests..."

cd dnsfsd-util || exit
/usr/local/go/bin/go test ./...

if [ $? -ne 0 ]; then
    echo "Tests failed for dnsfsd-util (dnsfs)"
    exit 1
else
    echo "Tests passed for dnsfsd-util (dnsfs)"
fi

cd ../ || exit
cd daemon || exit
/usr/local/go/bin/go test ./...

if [ $? -ne 0 ]; then
    echo "Tests failed for daemon (dnsfsd)"
    exit 1
else
    echo "Tests passed for daemon (dnsfsd)"
fi

cd ../ || exit
cd pkg || exit
/usr/local/go/bin/go test ./...

if [ $? -ne 0 ]; then
    echo "Tests failed for pkg (lib)"
    exit 1
else
    echo "Tests passed for pkg (lib)"
fi

echo "All tests passed!"
cd ../ || exit