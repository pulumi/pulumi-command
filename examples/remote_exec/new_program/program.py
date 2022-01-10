import os
import sys

data_dir = sys.argv[1]

TOTAL = 1
for fs_name in os.listdir(data_dir):
    # We don't use an f-string here to maintain python2 compatibility.
    fs_path = data_dir + "/" + fs_name
    with open(fs_path, 'r') as fs:
        data = fs.read()
        for number in data.splitlines(False):
            TOTAL *= int(number)

print(TOTAL)
