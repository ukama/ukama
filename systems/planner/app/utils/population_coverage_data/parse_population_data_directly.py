from osgeo import gdal, osr
from sqlalchemy import create_engine
from population_data_schema import PopulationData, FilesStatus
from sqlalchemy.orm import sessionmaker
import numpy as np

import glob

populationDataFilesPath = "./asdtest_geotiff_pop/*.tif" # change the path to the folder containing all the population geoTiff files.
SQLALCHEMY_DATABASE_URL = "mysql+mysqlconnector://root:abs-123@localhost:3308/planner_tool" # Change mySQL pass and db name which is planner tool
engine = create_engine(SQLALCHEMY_DATABASE_URL)
Session = sessionmaker(bind=engine)

def read_dir_for_geotiff_files():
    geotiff_files = glob.glob(populationDataFilesPath)
    return geotiff_files

def add_files_if_not_read(session, geotiff_files):
    for file in geotiff_files:
        existing_file = session.query(FilesStatus).filter_by(file_path=file).first()

        if existing_file is None:
            new_file = FilesStatus(file_path=file, parsed=False)
            session.add(new_file)
    session.commit()

def update_file_to_read(session, geotiff_file):
    existing_file = session.query(FilesStatus).filter_by(file_path=geotiff_file).first()

    if existing_file:
        existing_file.parsed = True
        session.commit()

def get_unparsed_files(session):
    unparsed_files = session.query(FilesStatus.file_path).filter_by(parsed=False).all()
    unparsed_files = [file[0] for file in unparsed_files]
    return unparsed_files

def store_parse_population_data(session, unparsed_data_files):
    box_data_list = []

    for image_path in unparsed_data_files:
        try:
            ds = gdal.Open(image_path)
            print(image_path)
            west, pixel_width, _, north, _, pixel_height = ds.GetGeoTransform()
            height, width = ds.RasterYSize, ds.RasterXSize

            band_data = ds.GetRasterBand(1).ReadAsArray()
            custom_value = 0
            band_data[np.isnan(band_data)] = custom_value

            box_data_list = []

            for row in range(height):
                for col in range(width):
                    longitude = west + col * pixel_width
                    latitude = north + row * pixel_height
                    population = band_data[row, col]
                    if population is not None and population > 0:
                        box_data = PopulationData(longitude=longitude, latitude=latitude, value=population)
                        box_data_list.append(box_data)

            ds.Close()
            update_file_to_read(session, image_path)
        except Exception as ex:
            print(ex)

    session.bulk_save_objects(box_data_list)
    session.commit()

    return {'message': 'Box data stored in the database'}

# Create session
session = Session()

# Read GeoTIFF files
geo_tiff_files = read_dir_for_geotiff_files()

# Add files to the database if not already read
add_files_if_not_read(session, geo_tiff_files)

# Get unparsed files
unparsed_files = get_unparsed_files(session)

print (unparsed_files)

store_parse_population_data(session, unparsed_files)