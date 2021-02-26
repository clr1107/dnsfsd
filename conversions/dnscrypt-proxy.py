# Convert lists for dnscrypt-proxy to dnsfsd. Only intended for blacklists.
# Run this python script with a dnscrypt-proxy ruleset file through stdin, it
# will then send to stdout the conversion.
# For example, on linux:
# `python3 dnscrypt-proxy.py < original.txt > converted.txt`

import sys
import os


def skip_line(line):
    return len(line) == 0 or line.isspace() or line[0] == '#'


def convert(line):
    return 'e;;' + line


def do():
    for line in sys.stdin:
        line = line.rstrip('\n').strip()

        if skip_line(line):
            continue

        sys.stdout.write(convert(line) + '\n')


if __name__ == '__main__':
    if os.isatty(sys.stdin.fileno()):
        print('no data piped in to script')
        exit(1)

    header = '''# Converted from dnscrypt-proxy to dnsfsd
# Script usage: python3 dnscrypt-proxy.py < original.txt > converted.txt'''

    print(header)
    do()
