# Convert lists for adblock plus 2.0 to dnsfsd. Only intended for blacklists.
# This script works with adblock exact addresses & domain name rules.
# Run this python script with an adblock ruleset file through stdin, it will
# then send to stdout the conversion.
# For example, on linux:
# `python3 adblock.py < original.txt > converted.txt`

import sys
import os


def skip_line(line):
    if len(line) == 0 or line.isspace():
        return True

    return not ((line[0] == '|' and line[len(line) - 1] == '|') or (
            line[0:2] == '||' and line[len(line) - 1] == '^'))


def convert(line):
    if line[0:2] == '||':
        return 'e;;' + line[2:len(line) - 1]

    return 'e;;' + line[1:len(line) - 1]


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

    header = '''# Converted from adblock plus 2.0 to dnsfsd
# Script usage: python3 adblock.py < original.txt > converted.txt'''

    print(header)
    do()
