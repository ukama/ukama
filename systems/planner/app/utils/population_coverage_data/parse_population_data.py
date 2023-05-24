from osgeo import gdal, osr
from sqlalchemy import create_engine
from population_data_schema import PopulationData, FilesStatus
from sqlalchemy.orm import sessionmaker
import glob

populationDataFilesPath = "E:/Projects/Freelance/UKAMA/test_geotiff_pop/*.tif"
SQLALCHEMY_DATABASE_URL = "mysql+mysqlconnector://root:MyNewPass@localhost/planner_tool"
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
        dataset = gdal.Open(image_path)

        # Get the image's spatial reference system (SRS) information
        srs = osr.SpatialReference()
        srs.ImportFromWkt(dataset.GetProjection())

        # Get the maximum and minimum coordinates from the image
        gt = dataset.GetGeoTransform()
        cols = dataset.RasterXSize
        rows = dataset.RasterYSize
        band = dataset.GetRasterBand(1)

        # Define the desired 30 by 30 meter square dimensions
        square_size = 30  # in meters

        pixel_size_x = gt[1]  # pixel width in x direction
        pixel_size_y = gt[5]  # pixel height in y direction (negative to indicate upwards direction)

        # Calculate the maximum and minimum coordinates
        max_x = gt[0] + (cols * gt[1])  # Maximum X coordinate
        min_x = gt[0]  # Minimum X coordinate
        max_y = gt[3]  # Maximum Y coordinate
        min_y = gt[3] + (rows * gt[5])  # Minimum Y coordinate

        # Calculate the number of pixels required to cover the square area
        num_pixels_x = int(square_size / pixel_size_x)
        num_pixels_y = int(square_size / pixel_size_y)
        import pdb; pdb.set_trace()
        # Calculate the number of boxes in X and Y directions
        num_boxes_x = int((max_x - min_x) / num_pixels_x)
        num_boxes_y = int((max_y - min_y) / num_pixels_y)

        # Calculate the width and height of each box
        box_width = (max_x - min_x) / num_boxes_x
        box_height = (max_y - min_y) / num_boxes_y

        # Calculate the coordinates of the boxes
        box_coordinates = []
        for i in range(num_boxes_x):
            for j in range(num_boxes_y):
                box_min_x = min_x + (i * box_width)
                box_max_x = box_min_x + box_width
                box_min_y = min_y + (j * box_height)
                box_max_y = box_min_y + box_height

                box_min_x_pixel = int((box_min_y - box_max_x) / gt[1])
                box_max_x_pixel = int((box_max_y - box_max_x) / gt[1])
                box_min_y_pixel = int((box_max_x - box_max_y) / gt[5])
                box_max_y_pixel = int((box_min_x - box_max_y) / gt[5])
                import pdb; pdb.set_trace()
                data = band.ReadAsArray(box_min_x_pixel, box_min_y_pixel, box_max_x_pixel - box_min_x_pixel, box_max_y_pixel - box_min_y_pixel)
                box_coordinates.append((box_max_x, box_min_x, box_max_y, box_min_y, data))

        # Print the coordinates of the boxes
        for i, box_coord in enumerate(box_coordinates):
            print(f"Box {i+1} coordinates (max_lon, min_lon, max_lat, min_lat): {box_coord}")

        # Store the box data in the database
        for i, box_coord in enumerate(box_coordinates):
            max_lon, min_lon, max_lat, min_lat, box_value = box_coord
            box_data = PopulationData(max_lon=max_lon, min_lon=min_lon, max_lat=max_lat, min_lat=min_lat, value=box_value)
            session.add(box_data)
        session.commit()
        session.close()

    return {'message': 'Box data stored in the database'}

geo_tiff_files = read_dir_for_geotiff_files()
add_files_if_not_read(geo_tiff_files)

unparsed_files = get_unparsed_files()

store_parse_population_data(unparsed_files)