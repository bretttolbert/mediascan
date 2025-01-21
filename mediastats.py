#!/usr/bin/python3
from __future__ import annotations
import yaml
import sys
import logging
import os
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
from typing import Dict, List, Set, Tuple

from dataclasses import dataclass

from dataclass_wizard import YAMLWizard  # type: ignore

"""
mediastats.py
Reads yaml file output by mediascan and calculates statistics
"""

log_level = logging.INFO
log_format = "[%(asctime)s] [%(levelname)s] %(name)s - %(message)s"
logging.basicConfig(
    level=logging.INFO,
    format=log_format,
    handlers=[
        logging.FileHandler("timebox.log", encoding="utf-8"),
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


def get_year_counts(data: Data) -> Dict[int, int]:
    counts: Dict[int, int] = {}
    for file in data.mediafiles:
        year = int(file.year)
        if year in counts:
            counts[year] += 1
        else:
            counts[year] = 1
    return counts


def get_genre_counts(
    data: Data, min_val: int, reverse_sort: bool = False
) -> Dict[str, int]:
    counts: Dict[str, int] = {}
    for file in data.mediafiles:
        genre = file.genre
        if genre in counts:
            counts[genre] += 1
        else:
            print(f"adding new genre '{genre}'")
            counts[genre] = 1
    counts_above_min_val = {k: v for (k, v) in counts.items() if v >= min_val}
    return dict(
        sorted(
            counts_above_min_val.items(), key=lambda item: item[1], reverse=reverse_sort
        )
    )


def get_file_sizes(data: Data) -> List[int]:
    return [int(f.size) for f in data.mediafiles]


def get_file_durations(data: Data) -> List[int]:
    return [int(f.duration) for f in data.mediafiles]


def get_years(data: Data):
    return [int(f.year) for f in data.mediafiles]


def plt_year_counts(data: Data):
    counts: Dict[int, int] = get_year_counts(data)
    df = pd.DataFrame({"year": counts.keys(), "count": counts.values()})
    df.plot.bar(x="year", y="count", rot=90)  # type: ignore


def plt_genre_counts(data: Data, min_val: int):
    counts = get_genre_counts(data, min_val)
    df = pd.DataFrame({"genre": counts.keys(), "count": counts.values()})
    df.plot.barh(x="genre", y="count", rot=0)  # type: ignore
    plt.yticks(fontsize=9)  # type: ignore


def plt_file_sizes(data: Data):
    sizes = get_file_sizes(data)
    plt.hist(sizes, bins=1000)  # type: ignore
    MEGABYTE = 10**6
    # MEBIBYTE = 2**20
    conversion_unit = MEGABYTE
    current_values = plt.gca().get_xticks()  # type: ignore
    plt.gca().set_xticklabels(  # type: ignore
        ["{:.2f}".format(x / conversion_unit) for x in current_values]
    )


def plt_file_durations(data: Data):
    durations = get_file_durations(data)
    plt.hist(durations, bins=range(800), color="b", edgecolor="red")  # type: ignore
    plt.xticks(rotation="vertical")  # type: ignore
    plt.xlim(0, 800)  # type: ignore
    plt.ylim(0, 100)  # type: ignore
    xtick_min = min(durations)
    # xtick_max = max(sizes)+1
    xtick_max = 800
    plt.xticks(np.arange(xtick_min, xtick_max, 10.0))  # type: ignore


def plt_year_vs_duration(data: Data):
    years = get_years(data)
    durations = get_file_durations(data)
    plt.scatter(x=years, y=durations)  # type: ignore
    plt.xlim(1950, 2022)  # type: ignore
    plt.ylim(0, 3000)  # type: ignore


def print_genres(data: Data):
    genres: Set[str] = set()  # type: ignore
    for f in data.mediafiles:
        genres.add(f.genre)
    print(sorted(genres))


def print_genre_counts(data: Data):
    counts = get_genre_counts(data, 1, reverse_sort=True)
    print(counts)


def get_album_paths(data: Data) -> Set[Tuple[str, str]]:
    """returns set of tuples of (album,path)"""
    albums: Set[Tuple[str, str]] = set()
    for f in data.mediafiles:
        albums.add((f.album, os.path.dirname(f.path)))
    return albums


def print_covers_by_size(data: Data):
    """Use this to find low-resolution album cover files that need to be updated to high-res"""
    album_paths = get_album_paths(data)
    cover_paths: List[Tuple[str, int]] = []
    for _, path in album_paths:
        cover_path = os.path.join(path, "cover.jpg")
        cover_size = os.path.getsize(cover_path)
        cover_paths.append((cover_path, cover_size))
    cover_paths = sorted(cover_paths, key=lambda item: item[1], reverse=True)
    for cover_path, size in cover_paths:
        print(f"{cover_path} {size}")


def main():
    if len(sys.argv) != 2:
        print("Usage: {0} <files yaml file>".format(sys.argv[0]))
    else:
        data = load_yaml_file(sys.argv[1])
        # print_covers_by_size(data)
        # print_genre_counts(data)
        # print_genres(data)
        # plt_year_counts(data)
        plt_genre_counts(data, 20)
        # plt_file_sizes(data)
        # plt_file_durations(data)
        # plt_year_vs_duration(data)
        plt.show()  # type: ignore


if __name__ == "__main__":
    main()
