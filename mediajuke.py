#!/usr/bin/python3
from __future__ import annotations
import logging
import random
import os
import yaml
import sys

from dataclasses import dataclass

from dataclass_wizard import YAMLWizard  # type: ignore

"""
mediajuke.py
Reads yaml file output by mediascan and plays a random file
"""

log_level = logging.INFO
log_format = "[%(asctime)s] [%(levelname)s] %(name)s - %(message)s"
logging.basicConfig(
    level=logging.INFO,
    format=log_format,
    handlers=[
        logging.FileHandler("mediajuke.log", encoding="utf-8"),
        logging.StreamHandler(),
    ],
)
log = logging.getLogger(__name__)


@dataclass
class Data(YAMLWizard):
    """
    Data dataclass

    """

    mediafiles: list[Mediafile]


@dataclass
class Mediafile:
    """
    Mediafile dataclass

    """

    path: str
    size: int
    format: str
    title: str
    artist: str
    albumartist: str
    album: str
    genre: str
    year: int
    duration: int


def load_yaml_file(yaml_fname: str) -> Data:
    data = None
    with open(yaml_fname, "r") as stream:
        try:
            data = Data.from_yaml(stream)  # type: ignore
        except yaml.YAMLError as exc:
            log.error(exc)
            sys.exit(1)
    return data  # type: ignore


def main():
    if len(sys.argv) != 3:
        print("Usage: {0} <player cmd> <files yaml file>".format(sys.argv[0]))
    else:
        player_cmd = sys.argv[1]
        file_yaml_path = sys.argv[2]
        data = load_yaml_file(file_yaml_path)
        file = random.choice(data.mediafiles)
        os.system(f'{player_cmd} "{file.path}"')


if __name__ == "__main__":
    main()
