
import os
from osgeo import gdal
from dotenv import load_dotenv
from typing import List

from app.elevation.schemas.elevation import Site, ElevationResponseSchema
from core import config


load_dotenv()

class SitesElevation:
    def __init__(self):
        self.HGT_FILES_PATH = os.getenv("HGT_DIR", config.get_config().HGT_FILES_PATH) + "/"
        self.HGT_EXTENSION = ".hgt"

    def get_elevations(self, sites: List[Site]) -> List[ElevationResponseSchema]:
        elevationResponses = []
        for site in sites:
            longitude = site['longitude']
            latitude = site['latitude']
            elevation = self.get_elevation_from_lon_lat(longitude, latitude)
            elevationResponses.append(
                ElevationResponseSchema(longitude=longitude, latitude=latitude, elevation=elevation)
            )
        return elevationResponses

    def get_elevation_from_lon_lat(self, longitude, latitude):
        try:
            if not (-180 <= longitude <= 180 and -90 <= latitude <= 90): # check for limit
                return 0
            
            ds = gdal.Open(self.get_raster_path(longitude, latitude), gdal.GA_ReadOnly)
            west, pixel_width, _, north, _, pixel_height = ds.GetGeoTransform()
            col = abs(int((longitude - west) / pixel_width))
            row = abs(int((north - latitude) / -pixel_height))
            band = ds.GetRasterBand(1)
            pRBB = band.ReadRaster(col, row, 1, 1, buf_type=gdal.GDT_Int32)
            pRBB = max(pRBB[0], 0)
            band = None
            ds = None
            return pRBB
        except Exception:
            return 0

    def get_raster_path(self, longitude, latitude):
        lon = str(longitude).strip().split('.')[0]
        lat = str(latitude).strip().split('.')[0]

        if lat.startswith("-"):
            path = "S" + str(int(lat.replace("-", "")) + 1).zfill(2)
        else:
            path = "N" + str(int(lat)).zfill(2)
        if lon.startswith("-"):
            path += "W" + str(int(lon.replace("-", "")) + 1).zfill(3) + self.HGT_EXTENSION
        else:
            path += "E" + str(int(lon)).zfill(3) + self.HGT_EXTENSION
        return self.HGT_FILES_PATH + path
