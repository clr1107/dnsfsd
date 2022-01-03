# DNSFSD
A DNS server that filters domain names by complex rule matching and forwards those that pass to another DNS server and those that don't get ignored ('sinkholed'). So this is essentially like a PiHole except it runs on a local system.

The server also has a configurable DNS cache.

Note: lots of commands here will require root permission.

### Install
To install this project it must be built using a Go compiler. First clone the repository: `git clone https://github.com/clr1107/dnsfsd.git`. Then run the `install.sh` script to build & move the binaries to an appropriate location. Lastly, run `dnsfs setup`. This has been tested on Ubuntu and will, for now, only work on Linux systems.

Once `dnsfs setup` has been run all directories and files (default configuration and systemd service file) will be created. From there one can start the server by using `dnsfsd` or `systemctl start dnsfsd`.

To start the server evert time the computer starts use `systemctl enable dnsfsd`

### Rules
Rule files, contained in `/etc/dnsfsd/rules`, follow a strict structure. Every new line is a new rule. So far there are three types of rules: regular expressions (`r`), contains (`c`), and equals (`e`). Lines that start with `#` are comments. The structure of a rule is as follows:
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

### Conversions
Rule files from other software can be converted to dnsfs using Python3 scripts located in the directory `conversions`
So far conversions for adblock dnscrypt-proxy, and hostfiles are done.

### dnsfs
The `dnsfs` command contains some useful utilities. 

#### dig
`dnsfs dig` which will allow one to test their rulesets by sending a fake (A type) DNS query.

#### download
`dnsfs download` will download an external (dnsfs) rule file and, with a given name, store it in `/etc/dnsfsd/rules/`. Note: the server will need to be restarted for the rule file to be loaded into a ruleset.

To download and convert a rule file from another piece of software one will have to use a external utility, such as `curl`, and use one of the conversion scripts and move it manually.

#### setup
`dnsfs setup` was discussed above.

#### clean
`dnsfs clean` deletes all logs that dnsfsd has created. Ensure dnsfsd is not running at the time.

#### log
`dnsfs log` outputs the log file, if it exists. It also supports `head` and `tail` functions. I.e. `dnsfs log 3` will read the first 3 lines only. `dnsfs log -- -3` will read the last 3 (in the standard order). Note the `--`, this is to signal that `-3` is a number and not a flag.

#### status
`dnsfs status` attempts to check with systemd if `dnsfsd` is running.