#/usr/bin/env/python

import yaml
import sys
import logging
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

log_level = logging.INFO
log_format = '[%(asctime)s] [%(levelname)s] %(name)s - %(message)s'
logging.basicConfig(
    level=logging.INFO,
    format=log_format,
    handlers=[
        logging.FileHandler("timebox.log", encoding="utf-8"),
        logging.StreamHandler()
    ]
)
log = logging.getLogger(__name__)

"""
mediastats.py
Reads yaml file output by mediascan and calculates statistics
"""

def load_yaml_file(yaml_fname):
    data = None
    with open(yaml_fname, 'r') as stream:
        try:
            data = yaml.safe_load(stream)
        except yaml.YAMLError as exc:
            log.error(exc)
            sys.exit(1)
    return data

def get_year_counts(data):
    counts = {}
    for file in data['mediafiles']:
        year = int(file['year'])
        if year in counts:
            counts[year] += 1
        else:
            counts[year] = 1
    return counts

def get_genre_counts(data):
    counts = {}
    for file in data['mediafiles']:
        genre = file['genre']
        if genre in counts:
            counts[genre] += 1
        else:
            counts[genre] = 1
    return dict(sorted(counts.items(), key=lambda item: item[1]))

def get_file_sizes(data):
    return [int(f['size']) for f in data['mediafiles']]

def get_file_durations(data):
    return [int(f['duration']) for f in data['mediafiles']]

def get_years(data):
    return [int(f['year']) for f in data['mediafiles']]

def plt_year_counts(data):
    counts = get_year_counts(data)
    df = pd.DataFrame({'year': counts.keys(), 'count': counts.values()})
    ax = df.plot.bar(x='year', y='count', rot='vertical')

def plt_genre_counts(data):
    counts = get_genre_counts(data)
    df = pd.DataFrame({'genre': counts.keys(), 'count': counts.values()})
    ax = df.plot.bar(x='genre', y='count', rot='vertical')

def plt_file_sizes(data):
    sizes = get_file_sizes(data)
    plt.hist(sizes, bins=1000)
    MEGABYTE = 10**6
    MEBIBYTE = 2**20
    conversion_unit = MEGABYTE
    current_values = plt.gca().get_xticks()
    plt.gca().set_xticklabels(["{:.2f}".format(x/conversion_unit) for x in current_values])

def plt_file_durations(data):
    durations = get_file_durations(data)
    plt.hist(durations, bins=range(800), color='b', edgecolor='red')
    plt.xticks(rotation='vertical')
    plt.xlim(0, 800)
    plt.ylim(0, 100)
    xtick_min = min(durations)
    #xtick_max = max(sizes)+1
    xtick_max = 800
    plt.xticks(np.arange(xtick_min, xtick_max, 10.0))

def plt_year_vs_duration(data):
    years = get_years(data)
    durations = get_file_durations(data)
    plt.scatter(x=years, y=durations)
    plt.xlim(1950, 2022)
    plt.ylim(0, 3000)

def main():
    data = load_yaml_file('files.yaml')
    plt_year_counts(data)
    #plt_genre_counts(data)
    #plt_file_sizes(data)
    #plt_file_durations(data)
    #plt_year_vs_duration(data)
    plt.show()

if __name__ == "__main__":
    main()
