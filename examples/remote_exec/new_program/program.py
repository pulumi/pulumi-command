#!/usr/bin/env python3

import sys
import os

data_dir = sys.argv[1]

total = 1
for fs_name in os.listdir(data_dir):
    fs_path = data_dir+"/"+fs_name
    with open(fs_path, 'r') as fs:
        data = fs.read()
        for number in data.splitlines(False):
            total *= int(number)

print(total)
