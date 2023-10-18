import numpy as np

import os.path
from itertools import product
import rasterio as rio
from rasterio import windows

infile = 'E:/Projects/Freelance/UKAMA/test_geotiff_pop/population_AF01_2018-10-01.tif' # change the path to the population geoTiff file to test.
out_path = './tmp/tiles'
output_filename = 'tile_{}-{}.tif'

def read_geotiff_replace_nan(window_data):
    # Replace NaN values with 0
    window_data[np.isnan(window_data)] = 0

    return window_data


def get_tiles(ds, width=256, height=256, map_units=False):

    if map_units:
        # Get pixel size
        px, py = ds.transform.a, -ds.transform.e
        width, height = int(width / px + 0.5) , int(height / px + 0.5)

    ncols, nrows = ds.meta['width'], ds.meta['height']

    offsets = product(range(0, ncols, width), range(0, nrows, height))
    big_window = windows.Window(col_off=0, row_off=0, width=ncols, height=nrows)
    for col_off, row_off in  offsets:
        window =windows.Window(col_off=col_off, row_off=row_off, width=width, height=height).intersection(big_window)
        transform = windows.transform(window, ds.transform)
        yield window, transform


with rio.open(infile) as inds:
    tile_width, tile_height = 10, 10  # in meters

    meta = inds.meta.copy()

    for window, transform in get_tiles(inds, tile_width, tile_height, map_units=True):

        meta['transform'] = transform
        meta['width'], meta['height'] = window.width, window.height
        outpath = os.path.join(out_path,output_filename.format(int(window.col_off), int(window.row_off)))
        with rio.open(outpath, 'w', **meta) as outds:
            window_data = inds.read(window=window)
            window_data = read_geotiff_replace_nan(window_data)
            if window_data.max() > 0:
                outds.write(window_data)
