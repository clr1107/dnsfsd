#!/bin/bash

if [ $EUID -ne 0 ]; then
    echo "You must have root permissions to run this script"
    exit 1
fi

echo "Building binaries..."

cd dnsfsd-util || exit
/usr/local/go/bin/go build -o dnsfs .

if [ $? -ne 0 ]; then
    echo "Could not build dnsfsd-util"
    exit 1
fi
cd ../

cd daemon || exit
/usr/local/go/bin/go build -o dnsfsd .

if [ $? -ne 0 ]; then
    echo "Could not build daemon"
    exit 1
fi
cd ../

echo "Built binaries"

mv dnsfsd-util/dnsfs /usr/local/bin/ > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "Could not move the built dnsfs binary to /usr/local/bin"
    exit 1
fi

mv daemon/dnsfsd /usr/local/bin/ > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "Could not move the built dnsfsd binary to /usr/local/bin"
    exit 1
fi

echo "Moved binaries"
echo "Install finished! Run \`dnsfs setup\` before starting dnsfsd"