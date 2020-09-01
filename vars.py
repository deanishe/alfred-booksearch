#!/usr/bin/env python2
# encoding: utf-8
#
# Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
# MIT Licence. See http://opensource.org/licenses/MIT
#
# Created on 2020-05-26
#

"""Remove/add custom user variables before/after committing."""

from __future__ import print_function, absolute_import

import argparse
import json
import os
import plistlib
from subprocess import check_call
import sys


INFO_PLIST = os.path.join(os.path.dirname(__file__), 'info.plist')
VAR_CACHE = os.path.join(os.path.dirname(__file__), 'vars.json')

WHITELIST = ('ACTION_DEFAULT', 'ACTION_ALT')
DELETE_PREFIXES = ('ACTION_',)
CLEAR_PREFIXES = ('USER_',)


def log(s, *args, **kwargs):
    """Log to STDOUT."""
    if args:
        s = s % args
    elif kwargs:
        s = s % kwargs

    print(s, file=sys.stdout)


def save_vars():
    """Save variables from info.plist to vars.json."""
    data = plistlib.readPlist(INFO_PLIST)
    var = {}
    for key, value in data['variables'].items():
        if not value:  # ignore empty variables
            continue

        if key in WHITELIST:
            continue

        for prefix in DELETE_PREFIXES:
            if key.startswith(prefix):
                var[key] = value
                del data['variables'][key]
                log('deleted %s', key)

        for prefix in CLEAR_PREFIXES:
            if key.startswith(prefix):
                var[key] = value
                data['variables'][key] = ''
                log('cleared %s', key)

    if var:
        with open(VAR_CACHE, 'wb') as fp:
            json.dump(var, fp, indent=2, sort_keys=True, separators=(',', ': '))

        plistlib.writePlist(data, INFO_PLIST)


def add_vars():
    """Save variables from info.plist to vars.json."""
    with open(VAR_CACHE) as fp:
        var = json.load(fp)
    data = plistlib.readPlist(INFO_PLIST)
    for key, value in var.items():
        data['variables'][key] = value
        log('set %s', key)

    plistlib.writePlist(data, INFO_PLIST)


def reload():
    """Tell Alfred to reload workflow."""
    s = """
        tell application id "com.runningwithcrayons.Alfred"
            reload workflow "net.deanishe.alfred.goodreads"
        end tell
        """
    check_call(['/usr/bin/osascript', '-e', s])


def parse_args():
    """Handle CLI arguments."""
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('-s', '--save', action='store_true',
                        help='save variables and remove them from info.plist')
    parser.add_argument('-a', '--add', action='store_true',
                        help='re-add variables to info.plist')
    parser.add_argument('-r', '--reload', action='store_true',
                        help='tell Alfred to reload workflow')
    return parser.parse_args()


def main():
    """Run script."""
    args = parse_args()
    if args.save:
        save_vars()
    else:
        add_vars()

    if args.reload:
        reload()


if __name__ == '__main__':
    main()
