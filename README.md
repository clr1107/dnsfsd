# DNSFSD
A DNS server that filters domain names by pattern matching, with regular expressions, and forwards those that pass to another DNS server, and those that don't get ignored. So this is essentially like a PiHole except it runs on a local system.

### Install
To install this project it must be built using a Go compiler. Run `install.sh` to build & move the binaries to an appropriate location; then run `dnsfs setup`.