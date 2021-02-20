#!/bin/bash

if [ "$EUID" -ne 0 ]; then
    echo "Run this script as root: sudo bash install.sh"
    exit 1
fi

echo "Creating necessary directories"
mkdir -p /etc/dnsfd/patterns

if [ $? -ne 0 ]; then
    echo "Could not create directories /etc/dnsfsd/..."
    exit 1
fi

mkdir -p /var/log/dnsfsd
if [ $? -ne 0 ]; then
    echo "Could not create directory /var/log/dnsfsd"
    exit 1
fi

echo "Moving binary & updating permissions"
mv dnsfsd /bin/dnsfsd
chmod +x /bin/dnsfsd

if [ $? -ne 0 ]; then
    echo "Could not move `dnsfsd` to `/bin/dnsfsd`"
    exit 1
fi
