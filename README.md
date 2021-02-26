# DNSFSD
A DNS server that filters domain names by rule matching, and forwards those that pass to another DNS server, and those that don't get ignored. So this is essentially like a PiHole except it runs on a local system.

### Install
To install this project it must be built using a Go compiler. Run `install.sh` to build & move the binaries to an appropriate location; then run `dnsfs setup`.

### Rules
Rule files, contained in `/etc/dnsfsd/rules`, follow a strict structure. Every new line is a new rule. So far there are three types of rules: regular expressions (`r`), contains (`c`), and equal (`e`). The structure of a rule is as follows:
```
t;w;<rule here>
```
where `t` is the opcode, and `w` is optional and signals that it's a whitelist rule rather than a blacklist. **All whitelists take precedence over blacklists**

E.g.
```
r;;[0-9]\.google\..*
e;w;456.google.com
```
The first rule would blacklist all domains following that regular expression pattern and the second rule would whitelist the domain `456.google.com`. Rules are case-insensitive.

Note that the whitelist signal is blank in the first, this is equal to the following expressions: `r;[0-9]\.google\..*` and `r;X;[0-9]\.google\..*` where X is any string, as if it is not `w` (or not present) it is simply ignored and interpreted as a blacklist signal.

### Rules
Rulesets from other software can be converted to dnsfsd using Python3.x scripts located in the directory `conversions`
So far conversions for adblock & dnscrypt-proxy are done. A list from github.com/notracking/hosts-blocklists has been converted (the dnscrypt-proxy one) and is in the directory `lists`