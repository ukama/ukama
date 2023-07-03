from osgeo import gdal, osr
from sqlalchemy import create_engine
from population_data_schema import PopulationData, FilesStatus
from sqlalchemy.orm import sessionmaker
import numpy as np
import csv
import concurrent.futures
import glob
import os

populationDataFilesPath = "E:/Projects/Freelance/UKAMA/test_geotiff_pop/*.tif" # change the path to the folder containing all the population geoTiff files.
SQLALCHEMY_DATABASE_URL = "mysql+mysqlconnector://root:MyNewPass@localhost/planner_tool" # Change mySQL pass and db name which is planner tool
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


def store_parse_population_data(session, unparsed_data_files, start_id):
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
                        box_data = (longitude, latitude, population)
                        box_data_list.append(box_data)

            
            start_id = save_to_csv(box_data_list, image_path, start_id)
            update_file_to_read(session, image_path)
        except Exception as ex:
            print(ex)
        finally:
            box_data_list = None
            ds = None
    
def save_to_csv(box_data_list, csv_filename, start_id):
    
    file_name = os.path.join(os.path.dirname(csv_filename), os.path.splitext(os.path.basename(csv_filename))[0])+".csv"
    with open(file_name, 'w', newline='') as csv_file:
        writer = csv.writer(csv_file, delimiter=';', quoting=csv.QUOTE_MINIMAL)
        writer.writerow(["id", "longitude", "latitude", "value"])
        for idx, (longitude, latitude, value) in enumerate(box_data_list, start=1):
            writer.writerow([start_id, longitude, latitude, value])
            start_id += 1
    
    return start_id
# Create session
session = Session()
start_id = 97759809
# Read GeoTIFF files
geo_tiff_files = read_dir_for_geotiff_files()

# Add files to the database if not already read
add_files_if_not_read(session, geo_tiff_files)

# Get unparsed files
unparsed_files = get_unparsed_files(session)
print (unparsed_files)


store_parse_population_data(session, unparsed_files, start_id)