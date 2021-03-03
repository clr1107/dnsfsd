# Convert blocking hostfiles to dnsfsd.
# Run this python script with a hostfile through stdin, it will then send to
# stdout the conversion.
# If any directive doesn't point to 0.0.0.0 or 127.0.0.1 or localhost then it
# will count as a whitelist instruction
# For example, on linux: `python3 hostfile.py < original.txt > converted.txt`

import sys
import os


def skip_line(line):
    return len(line) == 0 or line.isspace() or line[0] == '#'


def convert(domain, whitelist=False):
    s = 'e;'

    if whitelist:
        s += 'w'

    s += ';' + domain
    return s


def do():
    for line in sys.stdin:
        line = line.rstrip('\n').strip()

        if skip_line(line):
            continue

        parts = line.split(' ')

        if len(parts) != 2:
            continue

        if parts[0] == 'localhost' and parts[1] == '127.0.0.1':
            continue

        sys.stdout.write(convert(parts[1], parts[0] != '0.0.0.0') + '\n')


if __name__ == '__main__':
    if os.isatty(sys.stdin.fileno()):
        print('no data piped in to script')
        exit(1)

    header = '''# Converted hostfile to dnsfsd
# Script usage: python3 hostfile.py < original.txt > converted.txt'''

    print(header)
    do()
