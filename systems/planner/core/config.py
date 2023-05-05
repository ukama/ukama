import os
from pydantic import BaseSettings


class Config(BaseSettings):
    ENV: str = "development"
    DEBUG: bool = True
    APP_HOST: str = "0.0.0.0"
    APP_PORT: int = 8000
    RF_SERVER_PATH: str = "/home/ubuntu/Signal-Server"
    SOLAR_DATA_DIR: str = "/data/solarData"
    SDF_FILES_PATH: str = "/data/sdfData"
    HGT_FILES_PATH: str = "/data/hgtData"
    OUTPUT_PATH: str = "/home/ubuntu/output/"
    TEMP_FOLDER: str = "/tmp/planner/output/"


class DevelopmentConfig(Config): # Use this for docker 
    RF_SERVER_PATH: str = "/var/www/server/Signal-Server"
    SOLAR_DATA_DIR: str = "/data/solarData"
    SDF_FILES_PATH: str = "/data/sdfData"
    HGT_FILES_PATH: str = "/data/hgtData"
    OUTPUT_PATH: str = "/var/www/app/output/"
    TEMP_FOLDER: str = "/tmp/planner/output/"


class LocalConfig(Config): # use this for running locally
    RF_SERVER_PATH: str = "/home/ubuntu/Signal-Server"
    SOLAR_DATA_DIR: str = "E:/Projects/Freelance/UKAMA/Ukama-fork/systems/planner/data/solarData"
    SDF_FILES_PATH: str = "/data/sdfData"
    HGT_FILES_PATH: str = "/data/hgtData"
    OUTPUT_PATH: str = "/home/ubuntu/output/"
    TEMP_FOLDER: str = "./tmp/planner/output/"


class ProductionConfig(Config): # use this for production. 
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

