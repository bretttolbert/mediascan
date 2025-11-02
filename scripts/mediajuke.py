from __future__ import annotations
import logging
import random
import os
import sys

from mediascan import load_files_yaml

"""
mediajuke.py
Reads yaml file output by mediascan and plays a random file
"""


def main():
    if len(sys.argv) != 3:
        print("Usage: {0} <player cmd> <files yaml file>".format(sys.argv[0]))
    else:
        player_cmd = sys.argv[1]
        file_yaml_path = sys.argv[2]
        files = load_files_yaml(file_yaml_path)
        file = random.choice(files.files)
        os.system(f'{player_cmd} "{file.path}"')


if __name__ == "__main__":
    main()
