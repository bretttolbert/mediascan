import yaml

from mediascan.src.mediafiles import MediaFiles


def load_files_yaml(yaml_fname: str) -> MediaFiles:
    """
    raises: yaml.YAMLError
    """
    files = None
    with open(yaml_fname, "r") as stream:
        files = MediaFiles.from_yaml(stream)  # type: ignore
    return files  # type: ignore
