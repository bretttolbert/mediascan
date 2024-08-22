import json
import yaml
import sys

import dataclass_wizard.wizard_cli as cli


def load_yaml_file(yaml_fname: str):
    data = None
    with open(yaml_fname, "r") as stream:
        try:
            data = yaml.safe_load(stream)
        except yaml.YAMLError as exc:
            print(exc)
            sys.exit(1)
    return data


data = load_yaml_file("files.yaml")
print(cli.PyCodeGenerator(file_contents=json.dumps(data), experimental=True).py_code)
