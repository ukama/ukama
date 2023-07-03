import uuid
import subprocess
import os
from datetime import datetime
import numpy as np
from osgeo import gdal, osr
from typing import List
from dotenv import load_dotenv
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
import json
from coverage import CoverageResponseSchema
from population_data_schema import PopulationData
from coverage_enum import CoverageEnum

load_dotenv()

# {
#       "latitude": -1.287864,
#       "longitude": 36.822379,
#       "transmitter_height": 25
# }
# {
#     "north": -1.017864,
#     "east": 37.092379,
#     "west": 36.552379,
#     "south": -1.557864,
#     "url": "/home/ubuntu/output/output_2023-05-25_23-21-45/29dc5edf-45ad-4a72-8950-f294fda43523.png"
# }


# {
#       "latitude": -1.292499,
#       "longitude": 36.849288,
#       "transmitter_height": 25
# }
# {
#     "north": -1.022499,
#     "east": 37.119288,
#     "west": 36.579288,
#     "south": -1.562499,
#     "url": "/home/ubuntu/output/output_2023-05-25_23-22-33/5c9eb785-012e-4081-b6ee-d3d7bcf6a695.png"
# }

class SitesCoverage:
    def __init__(self):
        self.OUTPUT_PATH = "E:/Projects/Freelance/UKAMA/"
        self.TEMP_FOLDER = "./tmp/planner/output/"
        mysql_user_pass = "MyNewPass"
        mysql_user = "root"
        SQLALCHEMY_DATABASE_URL = f"mysql+mysqlconnector://{mysql_user}:{mysql_user_pass}@localhost/planner_tool" # Change mySQL pass and db name which is planner tool
        engine = create_engine(SQLALCHEMY_DATABASE_URL)
        Session = sessionmaker(bind=engine)
        self.SESSION = Session()

    def calculate_coverage(self) -> CoverageResponseSchema:
        try:
            output_folder_path = self.generate_output_folder()
            sites_coverage_list = []
                
            output_file_path1 = "E:/Projects/Freelance/UKAMA/TestingWithData/46e6b4cd-41f3-4392-bc94-fe59be3a259c"
            output_file_path2 = "E:/Projects/Freelance/UKAMA/TestingWithData/629c71ef-e3c0-4196-9b47-9a3077d129f7"


            sites_coverage_list.append(
                CoverageResponseSchema(
                    north=35.592822,
                    east=-2.622917,
                    south=35.052456,
                    west=-3.284583,
                    url=f"{output_file_path1}.png",
                    population_data= {}
                )
            )

            sites_coverage_list.append(
                CoverageResponseSchema(
                    north=35.586153,
                    east=-2.620139,
                    south=35.045791,
                    west=-3.281805,
                    url=f"{output_file_path2}.png",
                    population_data= {}
                )
            )
        except subprocess.CalledProcessError as e:
            raise e
        except Exception as ex:
            raise ex
        else:
            return self.merge_sites_output(
                sites_coverage_list, output_file_path1, output_folder_path, CoverageEnum.RECEIVE_POWER.value
            )
        
    def merge_sites_output(
        self, sites_coverage_list, output_file_path, output_folder_path, outputFunc
    ):
        """
        Merges multiple geo tiff files into a single output Tiff image file along with new coordinates.

        Args:
        - sites_coverage_list: A list of dictionaries containing the file url and coordinates for each input file.
        - output_file_path: The path where the intermediate geo tiff files were created and use it to get dcf path.
        - output_folder_path: The path where the output geo tiff file will be saved after merging.
        - outputFunc: A function to select minimum or maximum value depending on the mode of prediction.

        Returns:
        - output: The output from the `merge_geo_tiff_files` function, which will be CoverageResponseSchema.
        """
        responseArray = self.create_geo_tiff_files(sites_coverage_list)
        output_file_name = str(uuid.uuid4()) + "_merged.tif"
        merged_file_url = output_folder_path + output_file_name
        output = self.merge_geo_tiff_files(
            responseArray, merged_file_url, output_file_path, outputFunc
        )
        return output

    # region Using gdal to merge files.

    def translate(self, input_path, output_path, options):
        inputds = gdal.Open(input_path)
        gdal.Translate(output_path, inputds, options=options)

    def create_geo_tiff_files(self, inputsArr: List[CoverageResponseSchema]):
        responseImages = []

        for input in inputsArr:
            input_file_path = input.url
            output_file_name = self.get_image_name(
                input_file_path
            )  # name without extension
            output_file_url = self.OUTPUT_PATH + output_file_name + ".tif"
            west, south, east, north = input.west, input.south, input.east, input.north
            translate_options = gdal.TranslateOptions(
                format="GTiff",
                outputSRS="EPSG:4326",
                outputBounds=[west, north, east, south],
            )
            self.translate(input_file_path, output_file_url, options=translate_options)
            responseImages.append(
                {
                    "image_name": output_file_name,
                    "image_url": output_file_url,
                }
            )
        return responseImages

    def merge_geo_tiff_files(self, inputs, merged_tif_url, dcf_file_path, outputFunc):
        population_image_pixel_width, population_image_pixel_height = 0.0002777777777777777775, 0.0002777777777777777775 # this is default as given by the data downloaded for population

        dcf_file_url = dcf_file_path + ".dcf"
        color_map = self.load_dcf_file(dcf_file_url)
        output_folder_path = "E:/Projects/Freelance/UKAMA/"

        # Create a dictionary to store the combined values
        combined = {}
        interference_map = {}
        # Rounding of longitude and latitude to match with other images, otherwise no coordinate of any pixel were matching
        roundofCoord = 3
        output_data_interference = {}
        y_min, y_max, x_min, x_max = float("inf"), float("-inf"), float("inf"), float("-inf")
        for input_file in inputs:
            # Open the image using GDAL
            ds = gdal.Open(input_file["image_url"])
            image_name = input_file["image_name"]
            # Get the spatial extent of the image in geographic coordinates
            west, pixel_width, _, north, _, pixel_height = ds.GetGeoTransform()
            height, width = ds.RasterYSize, ds.RasterXSize
            
            # Calculate the longitude and latitude of each pixel
            data = np.empty((height, width, 4), dtype=np.uint8)
            for i in range(4):
                data[:, :, i] = ds.GetRasterBand(i + 1).ReadAsArray()

            for row in range(height):
                for col in range(width):
                    # Calculate the longitude and latitude of the center of the pixel
                    longitude = west + col * pixel_width    
                    latitude = north + row * pixel_height
                    y_min, y_max = min(y_min, latitude), max(y_max, latitude)
                    x_min, x_max = min(x_min, longitude), max(x_max, longitude)
                    longitude, latitude = round(longitude, roundofCoord), round(latitude, roundofCoord)
                    key = (0, longitude, latitude), (1, longitude, latitude), (2, longitude, latitude), (3, longitude, latitude)
                    if data[row, col, 3] != 0:
                        rgb = tuple(data[row, col, :3])
                        if key in combined:
                            if outputFunc == CoverageEnum.PATH_LOSS.value:
                                if color_map[rgb] < color_map[tuple(combined[key][0][:3])]:
                                    combined[key] = [ tuple(list(rgb) + [data[row, col, 3]]), image_name]
                            else:
                                if color_map[rgb] > color_map[tuple(combined[key][0][:3])]:
                                    combined[key] = [ tuple(list(rgb) + [data[row, col, 3]]), image_name]
                                if outputFunc == CoverageEnum.RECEIVE_POWER.value and color_map[rgb] == color_map[tuple(combined[key][0][:3])] and image_name != combined[key][1]:
                                    interference_map[key] = image_name + "_and_" + combined[key][1]
                        else:
                            combined[key] = [tuple(list(rgb) + [data[row, col, 3]]), image_name]
            ds = None

        nrows = int(abs((y_max - y_min) / pixel_height))
        ncols = int(abs((x_max - x_min) / pixel_width))
        output_data = np.zeros((nrows, ncols, 4), dtype=np.uint8)
        

        # Creating combined geotiff file for cloud rf prediction
        for row in range(nrows):
            for col in range(ncols):
                lon = round(x_min + col * pixel_width, roundofCoord)
                lat = round(y_max + row * pixel_height, roundofCoord)
                key = (0, lon, lat), (1, lon, lat), (2, lon, lat), (3, lon, lat)
                if key in combined:
                    output_data[row, col, :3] = combined[key][0][:3]
                    output_data[row, col, 3] = combined[key][0][3]
                    if key in interference_map:
                        if interference_map[key] not in output_data_interference:
                            output_data_interference[interference_map[key]] = np.full((nrows, ncols, 1), np.nan, dtype=np.float32)
                        output_data_interference[interference_map[key]][row, col, 0] = 5

        self.create_geo_tiff_file(merged_tif_url, 
                            ncols, 
                            nrows,
                            4,
                            gdal.GDT_Byte, 
                            pixel_width, 
                            pixel_height, 
                            x_min, 
                            y_max, 
                            output_data)

        # Creating interference geotiffs
        interference_files_urls = {}
        for interference_data_key in output_data_interference:
            output_url = output_folder_path + interference_data_key + ".tif"
            self.create_geo_tiff_file(output_url, 
                                      ncols, 
                                      nrows,
                                      1,
                                      gdal.GDT_Float32, 
                                      pixel_width, 
                                      pixel_height, 
                                      x_min, 
                                      y_max, 
                                      output_data_interference[interference_data_key])
            interference_files_urls[interference_data_key] = output_url


        # Creating population coverage geotiffs
        pop_data = self.get_population_data(x_min, x_max, y_min, y_max)
        
        pop_data_nrows = int(abs((y_max - y_min) / population_image_pixel_height))
        pop_data_ncols = int(abs((x_max - x_min) / population_image_pixel_width))
        population_output_data = {}
        population_output_value = {}
        population_output_total_boxes = {}
        for input_file in inputs:
            population_output_data[input_file["image_name"]] = np.full((pop_data_nrows, pop_data_ncols, 1), np.nan, dtype=np.float32)
            population_output_value[input_file["image_name"]] = 0
            population_output_total_boxes[input_file["image_name"]] = 0

        # Creating separate geotiff files for population prediction
        for point in pop_data:
            longitude = point.longitude
            latitude = point.latitude
            population = point.value

            # Convert longitude and latitude to pixel coordinates
            x = int((longitude - x_min) / population_image_pixel_width)
            y = int((latitude - y_min) / population_image_pixel_height)
            roundedLon = round(longitude, roundofCoord)
            roundedLat = round(latitude, roundofCoord)
            key = (0, roundedLon, roundedLat), (1, roundedLon, roundedLat), (2, roundedLon, roundedLat), (3, roundedLon, roundedLat)
            # Write the population value to the raster
            if key in combined:
                population_output_data[combined[key][1]][pop_data_nrows - y - 1, x, 0] = population
                population_output_value[combined[key][1]] += population
                population_output_total_boxes[combined[key][1]] += 1

        siteNumber = 1
        population_data_dic = {}
        for input_file in inputs:
            self.create_geo_tiff_file(input_file["image_url"], 
                                      pop_data_ncols, 
                                      pop_data_nrows,
                                      1,
                                      gdal.GDT_Float32, 
                                      population_image_pixel_width, 
                                      -population_image_pixel_height, 
                                      x_min, 
                                      y_max, 
                                      population_output_data[input_file["image_name"]])
            
            print("site "+ str(siteNumber) + " population data file url: ", input_file["image_url"])
            siteNumber = siteNumber + 1 
            population_data_dic[input_file["image_name"]] = { "url": input_file["image_url"], "population_covered": population_output_value[input_file["image_name"]], "total_30m_areas_covered": population_output_total_boxes[input_file["image_name"]] }

        print("East: ", x_max)
        print("West: ", x_min)
        print("South: ", y_min)
        print("North: ", y_max)
        print("Merged file url: ", merged_tif_url)
        print("population output: ", json.dumps(population_data_dic, indent = 4))
        print("interference output: ", json.dumps(interference_files_urls, indent = 4))

        # Clean up
        output_ds = None
        combined = None
        ds = None
        return {
            "east": x_max,
            "west": x_min,
            "south": y_min,
            "north": y_max,
            "url": merged_tif_url,
            "population_data": population_data_dic
        }

    # endregion

    # region Helper methods
    def create_geo_tiff_file(self, output_url, cols, rows, bands, type, pixel_width, pixel_height, x_min, y_max, data):
        driver = gdal.GetDriverByName("GTiff")
        output_ds = driver.Create(output_url, cols, rows, bands, type)
        srs = osr.SpatialReference()
        srs.ImportFromEPSG(4326)
        output_ds.SetProjection(srs.ExportToWkt())
        output_ds.SetGeoTransform((x_min, pixel_width, 0, y_max, 0, pixel_height))

        for i in range(bands):
            output_band = output_ds.GetRasterBand(i + 1)
            output_band.WriteArray(data[:, :, i])
        output_ds.FlushCache()
        output_ds = None
        output_band = None

    def load_dcf_file(self, url):
        # Load DCF file
        with open(url, "r") as f:
            lines = f.readlines()

        # Create a dictionary to map pixel values to RGB colors
        color_map = {}
        for line in lines:
            r, g, b = list(
                map(int, [i.strip() for i in line.split(":")[-1].split(",")])
            )
            value = int(line.strip().split(":")[0])
            color_map[(r, g, b)] = value
        return color_map

    def get_image_name(self, image_url):
        file_name_with_ext = os.path.basename(image_url)
        file_name = file_name_with_ext.split(".")[0]
        return file_name

    def generate_output_folder(self) -> str:
        now = datetime.now()
        date_time_str = now.strftime("%Y-%m-%d_%H-%M-%S")

        folder_name = f"output_{date_time_str}"
        output_folder_path = self.OUTPUT_PATH + folder_name + "/"
        subprocess.check_output("mkdir -p " + output_folder_path)
        subprocess.check_output("mkdir -p " + self.TEMP_FOLDER)
        return output_folder_path
    
    def remove_temp_folder(self):
        print("removing temporary content in: "+ self.TEMP_FOLDER)
        subprocess.check_output("rm -rf " + self.TEMP_FOLDER)
    
    def filter_coordinates(self, pop_data, longitude, latitude):
        filtered_coordinates = [
            pop_coord for pop_coord in pop_data if round(pop_coord.longitude, 3) == longitude and round(pop_coord.latitude, 3) == latitude
        ]
        if len(filtered_coordinates) >= 1:
            return filtered_coordinates[0].get('value')
        
        return None

    def get_population_data(self, west, east, south, north):
        data = self.SESSION.query(PopulationData).filter(
            PopulationData.longitude >= west,
            PopulationData.longitude <= east,
            PopulationData.latitude >= south,
            PopulationData.latitude <= north
        ).all()
        
        return data
    
# endregion

sitescov = SitesCoverage()
print(json.dumps(sitescov.calculate_coverage(), indent = 4))