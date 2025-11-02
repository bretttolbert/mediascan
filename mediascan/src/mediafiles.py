from dataclasses import dataclass
from dataclass_wizard import YAMLWizard  # type: ignore

from mediascan.src.mediafile import MediaFile


@dataclass
class MediaFiles(YAMLWizard):
    """
    MediaFiles dataclass
    Data model for entire files YAML file output by mediascan.go

    """

    mediafiles: list[MediaFile]
