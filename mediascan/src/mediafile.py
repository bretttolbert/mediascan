from dataclasses import dataclass


@dataclass
class MediaFile:
    """
    MediaFile dataclass

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
