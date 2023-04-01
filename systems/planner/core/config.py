import os
from pydantic import BaseSettings


class Config(BaseSettings):
    ENV: str = "development"
    DEBUG: bool = True
    APP_HOST: str = "0.0.0.0"
    APP_PORT: int = 8000
    RF_SERVER_PATH: str = "/home/ubuntu/Signal-Server"
    SDF_FILES_PATH: str = "/data/sdfData"
    OUTPUT_PATH: str = "/home/ubuntu/output/"
    TEMP_FOLDER: str = "/tmp/planner/output/"


class DevelopmentConfig(Config):
    RF_SERVER_PATH: str = "/var/www/app/Signal-Server"
    SDF_FILES_PATH: str = "/data/sdfData"
    OUTPUT_PATH: str = "/var/www/app/output/"
    TEMP_FOLDER: str = "/var/www/app/tmp/planner/output/"


class LocalConfig(Config): 
    RF_SERVER_PATH: str = "/home/ubuntu/Signal-Server"
    SDF_FILES_PATH: str = "/data/sdfData"
    OUTPUT_PATH: str = "/home/ubuntu/output/"
    TEMP_FOLDER: str = "./tmp/planner/output/"


class ProductionConfig(Config):
    DEBUG: str = False
    # TODO add according to requirement


def get_config() -> Config:
    env = os.getenv("ENV", "local")
    config_type = {
        "dev": DevelopmentConfig(),
        "local": LocalConfig(),
        "prod": ProductionConfig(),
    }
    return config_type[env]


config: Config = get_config()

