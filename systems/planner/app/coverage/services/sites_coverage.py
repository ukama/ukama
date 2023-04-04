import uuid
import subprocess
import os
from datetime import datetime
import numpy as np
from osgeo import gdal, osr
from typing import List
from dotenv import load_dotenv

from app.coverage.schemas.coverage import Site, CoverageResponseSchema
from app.coverage.enums.coverage import CoverageEnum
from core import config

load_dotenv()

class SitesCoverage:
    def __init__(self):
        self.RF_SERVER_PATH = config.get_config().RF_SERVER_PATH
        self.SDF_FILES_PATH = os.getenv("SDF_DIR", config.get_config().SDF_FILES_PATH)
        self.OUTPUT_PATH = config.get_config().OUTPUT_PATH
        self.TEMP_FOLDER = config.get_config().TEMP_FOLDER

    def calculate_coverage(self, mode, sites: list[Site]) -> CoverageResponseSchema:
        try:
            output_folder_path = self.generate_output_folder()
            sites_coverage_list = []
            for site in sites:
                params = ["-sdf", self.SDF_FILES_PATH]
                if site['latitude']:
                    params.extend(["-lat", str(site['latitude'])])
                if site['longitude']:
                    params.extend(["-lon", str(site['longitude'])])
                if site['transmitter_height']:
                    params.extend(["-txh", str(site['transmitter_height'])])
                params.extend(["-f", "900", "-rt", "-110"])
                if mode == CoverageEnum.PATH_LOSS.value:
                    params.extend(["-erp", "0"])
                    output_func = "min"
                elif mode == CoverageEnum.FIELD_STRENGTH.value:
                    params.extend(["-erp", "20"])
                    output_func = "max"
                else:
                    params.extend(["-erp", "20", "-dbm"])
                    output_func = "max"
                params.extend(["-m", "-R", "30", "-res", "1200", "-pm", "1"])

                output_file_name = str(uuid.uuid4())
                output_file_path = f"{output_folder_path}{output_file_name}"

                params.extend(["-o", output_file_path])

                rf_find_cov_command = f"{self.RF_SERVER_PATH}/src/signalserver"
                print(f"Running {rf_find_cov_command} {' '.join(params)}")
                result = subprocess.check_output(
                    [rf_find_cov_command] + params, stderr=subprocess.STDOUT, text=True
                )

                output = result.split("|")

                print("Converting to PNG...")
                rf_convert_command = f"convert {output_file_path}.ppm -transparent white -channel Alpha PNG32:{output_file_path}.png"
                print(f"Running {rf_convert_command}")
                subprocess.check_call(rf_convert_command, shell=True)
                print(output)
                if len(output[1:-1]) == 4:
                    sites_coverage_list.append(
                        CoverageResponseSchema(
                            north=output[1],
                            east=output[2],
                            south=output[3],
                            west=output[4],
                            url=f"{output_file_path}.png",
                        )
                    )
        except subprocess.CalledProcessError as e:
            raise e
        except Exception as ex:
            raise ex
        else:
            if len(sites_coverage_list) > 1:
                return self.merge_sites_output(
                    sites_coverage_list, output_file_path, output_folder_path, output_func
                )
            else:
                return sites_coverage_list[0]
        finally:
            self.remove_temp_folder()

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
        test = gdal.Translate(output_path, inputds, options=options)

        inputds = None
        if test == None:
            return False
        test = None

        return True

    def create_geo_tiff_files(self, inputsArr: List[CoverageResponseSchema]):
        responseImages = []

        for input in inputsArr:
            input_file_path = input.url
            output_file_name = self.get_image_name(
                input_file_path
            )  # name without extension
            output_file_url = self.TEMP_FOLDER + output_file_name + ".tif"
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
        dcf_file_url = dcf_file_path + ".dcf"
        color_map = self.load_dcf_file(dcf_file_url)

        # Create a dictionary to store the combined values
        combined = {}
        # Rounding of longitude and latitude to match with other images, otherwise no coordinate of any pixel were matching
        roundofCoord = 3

        y_min = float("inf")
        y_max = float("-inf")
        x_min = float("inf")
        x_max = float("-inf")
        # Loop through the input files
        for input_file in inputs:
            # Open the image using GDAL
            ds = gdal.Open(input_file["image_url"])

            # Get the spatial extent of the image in geographic coordinates
            west, pixel_width, _, north, _, pixel_height = ds.GetGeoTransform()
            height, width = ds.RasterYSize, ds.RasterXSize

            # Calculate the longitude and latitude of each pixel
            data_red = ds.GetRasterBand(1).ReadAsArray()
            data_green = ds.GetRasterBand(2).ReadAsArray()
            data_blue = ds.GetRasterBand(3).ReadAsArray()
            data_alpha = ds.GetRasterBand(4).ReadAsArray()

            for row in range(height):
                for col in range(width):
                    # Calculate the longitude and latitude of the center of the pixel
                    longitude = west + col * pixel_width
                    latitude = north + row * pixel_height
                    y_min = min(y_min, latitude)
                    y_max = max(y_max, latitude)
                    x_min = min(x_min, longitude)
                    x_max = max(x_max, longitude)
                    longitude = round(longitude, roundofCoord)
                    latitude = round(latitude, roundofCoord)
                    red = (0, longitude, latitude)
                    green = (1, longitude, latitude)
                    blue = (2, longitude, latitude)
                    alpha = (3, longitude, latitude)

                    if data_alpha[row, col] != 0:
                        rgb = (
                            data_red[row, col],
                            data_green[row, col],
                            data_blue[row, col],
                        )
                        if red in combined:
                            if outputFunc == "max":
                                if (
                                    color_map[rgb]
                                    > color_map[
                                        (combined[red], combined[green], combined[blue])
                                    ]
                                ):
                                    combined[red], combined[green], combined[blue] = rgb
                                    combined[alpha] = data_alpha[row, col]
                            else:
                                if (
                                    color_map[rgb]
                                    < color_map[
                                        (combined[red], combined[green], combined[blue])
                                    ]
                                ):
                                    combined[red], combined[green], combined[blue] = rgb
                                    combined[alpha] = data_alpha[row, col]

                        else:
                            combined[red], combined[green], combined[blue] = rgb
                            combined[alpha] = data_alpha[row, col]

        nrows = int(abs((y_max - y_min) / pixel_height))
        ncols = int(abs((x_max - x_min) / pixel_width))
        output_data = np.zeros((nrows, ncols, 4), dtype=np.uint8)

        for row in range(nrows):
            for col in range(ncols):
                lon = round(x_min + col * pixel_width, roundofCoord)
                lat = round(y_max + row * pixel_height, roundofCoord)
                red = (0, lon, lat)
                green = (1, lon, lat)
                blue = (2, lon, lat)
                alpha = (3, lon, lat)
                if red in combined:
                    output_data[row, col, 0] = int(combined[red])
                    output_data[row, col, 1] = int(combined[green])
                    output_data[row, col, 2] = int(combined[blue])
                    output_data[row, col, 3] = int(combined[alpha])

        driver = gdal.GetDriverByName("GTiff")
        output_ds = driver.Create(merged_tif_url, ncols, nrows, 4, gdal.GDT_Byte)
        srs = osr.SpatialReference()
        srs.ImportFromEPSG(4326)
        output_ds.SetProjection(srs.ExportToWkt())
        output_ds.SetGeoTransform((x_min, pixel_width, 0, y_max, 0, pixel_height))

        for i in range(4):
            output_band = output_ds.GetRasterBand(i + 1)
            output_band.WriteArray(output_data[:, :, i])
        output_ds.FlushCache()

        print("East: ", x_max)
        print("West: ", x_min)
        print("South: ", y_min)
        print("North: ", y_max)
        print("Merged file url: ", merged_tif_url)

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
        }

    # endregion

    # region Helper methods

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
        subprocess.check_output("mkdir -p " + output_folder_path, shell=True)
        subprocess.check_output("mkdir -p " + self.TEMP_FOLDER, shell=True)
        return output_folder_path
    
    def remove_temp_folder(self):
        print("removing temporary content in: "+ self.TEMP_FOLDER)
        subprocess.check_output("rm -rf " + self.TEMP_FOLDER, shell=True)

# endregion
