#!/usr/bin/env python3

import os
import sqlite3
import sys

here = os.path.abspath(os.path.dirname(__file__))


def restore(database, file):
    db = sqlite3.connect(database)
    with open(os.path.join(here, 'drop_all.sql')) as f:
        db.executescript(f.read())
    with open(file) as f:
        db.executescript(f.read())
    db.close()


def backup(database, file):
    db = sqlite3.connect(database)
    with open(file, 'w') as f:
        f.write('\n'.join(db.iterdump()))
    db.close()


def main():
    if sys.argv[1] == 'restore':
        restore(sys.argv[2], sys.argv[3])
    elif sys.argv[1] == 'backup':
        backup(sys.argv[2], sys.argv[3])
    else:
        os._exit(1)


if __name__ == '__main__':
    main()
