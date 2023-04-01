from osgeo import gdal
import simplekml

# Set input GeoTIFF file path
input_file = "./tmp/output/testing_1_merged.tif"

# Set output KMZ file path
output_file = "./tmp/output/output3.kmz"

# Open the GeoTIFF file
ds = gdal.Open(input_file)

# Get the GeoTransform and projection information
gt = ds.GetGeoTransform()
proj = ds.GetProjection()

# Create a simplekml KML object
kml = simplekml.Kml()

# Create a ground overlay from the GeoTIFF file
ground = kml.newgroundoverlay(name="overlay")
ground.icon.href = input_file
ground.latlonbox.north = gt[3]
ground.latlonbox.south = gt[3] + gt[5] * ds.RasterYSize
import pdb; pdb.set_trace()
ground.latlonbox.east = gt[0] + gt[1] * ds.RasterXSize
ground.latlonbox.west = gt[0]

# Save the KMZ file
kml.savekmz(output_file)
