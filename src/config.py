from typing import Any
import yaml


class ConfigReader:
    def __init__(self, config_file: str = "config.yml"):
        self.config_file = config_file
        self.config_data = self._read_config()

    def _read_config(self) -> None:
        try:
            with open(self.config_file, "r") as file:
                config_data = yaml.safe_load(file)
                print("Configuration file loaded successfully.")
                return config_data
        except FileNotFoundError:
            print(f"Error: Configuration file '{self.config_file}' not found.")
            return None
        except yaml.YAMLError as e:
            print(f"Error parsing YAML file: {e}")
            return None

    def get_value(self, key) -> Any:
        if self.config_data:
            return self.config_data.get(key)
        else:
            print("Error: Configuration data not loaded.")
            return None
