#!/bin/bash
set -eu

run_tests () {
    oldwd="${PWD}"
    cd "${1}" || exit 1
    go test ./...
    cd "${oldwd}"
}

echo "Running tests..."

if run_tests ./dnsfsd-util; then
    echo "Tests passed for dnsfsd-util (dnsfs)"
else
    echo "Tests failed for dnsfsd-util (dnsfs)"
    exit 1
fi

if run_tests ./daemon; then
    echo "Tests passed for daemon (dnsfsd)"
else
    echo "Tests failed for daemon (dnsfsd)"
    exit 1
fi

if run_tests ./pkg; then
    echo "Tests passed for pkg (lib)"
else
    echo "Tests failed for pkg (lib)"
    exit 1
fi

echo "All tests passed!"
