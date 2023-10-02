from osgeo import gdal, osr
from sqlalchemy import create_engine
from population_data_schema import PopulationDataSimple, FilesStatus
from sqlalchemy.orm import sessionmaker
import numpy as np

import glob

populationDataFilesPath = "./UKAMA/test_geotiff_pop/*.tif" # change the path to the folder containing all the population geoTiff files.
SQLALCHEMY_DATABASE_URL = "mysql+mysqlconnector://root:MyNewPass@localhost/planner_tool" # Change mySQL pass and db name which is planner tool
engine = create_engine(SQLALCHEMY_DATABASE_URL)
Session = sessionmaker(bind=engine)
session = Session()
def read_dir_for_geotiff_files():
    geotiff_files = glob.glob(populationDataFilesPath)
    return geotiff_files

def add_files_if_not_read(geotiff_files):
    for file in geotiff_files:
        # Check if the file path already exists in the database
        existing_file = session.query(FilesStatus).filter_by(file_path=file).first()

        if existing_file is None:
            # Insert the file path into the database
            new_file = FilesStatus(file_path=file, parsed=False)
            session.add(new_file)
    session.commit()

    # Close the session
    session.close()

def update_file_to_read(geotiff_file):
    # Check if the file path already exists in the database
    existing_file = session.query(FilesStatus).filter_by(file_path=geotiff_file).first()

    if existing_file:
        existing_file.parsed = True  # Set 'parsed' column to True
        session.commit()
    # Close the session
    session.close()

def get_unparsed_files():
    try:
        # Retrieve a list of all the files that are not parsed
        unparsed_files = session.query(FilesStatus.file_path).filter_by(parsed=False).all()
        unparsed_files = [file[0] for file in unparsed_files]
        return unparsed_files
    finally:
        # Close the session
        session.close()


def store_parse_population_data(unparsed_data_files):
    # Load the GeoTIFF image and extract the box data
    # Load the GeoTIFF image
    for image_path in unparsed_data_files:
        ds = gdal.Open(image_path)

        # Get the spatial extent of the image in geographic coordinates
        west, pixel_width, _, north, _, pixel_height = ds.GetGeoTransform()
        height, width = ds.RasterYSize, ds.RasterXSize

        data = np.empty((height, width, 1), dtype=np.uint8)
        band_data = ds.GetRasterBand(1).ReadAsArray()
        custom_value = 0  # Specify the desired value to replace NaN
        band_data[np.isnan(band_data)] = custom_value
        
        for row in range(height):
            for col in range(width):
                # Calculate the longitude and latitude of the center of the pixel
                longitude = west + col * pixel_width
                latitude = north + row * pixel_height
                population = band_data[row, col]
                if population != None and 0 < population:
                    box_data = PopulationDataSimple(longitude=longitude, latitude=latitude, value=population)
                    session.add(box_data)
        session.commit()
        session.close()
        update_file_to_read(image_path)
        

    return {'message': 'Box data stored in the database'}

geo_tiff_files = read_dir_for_geotiff_files()
add_files_if_not_read(geo_tiff_files)

unparsed_files = get_unparsed_files()

store_parse_population_data(unparsed_files)