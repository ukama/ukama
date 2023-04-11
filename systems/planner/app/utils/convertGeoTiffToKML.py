from osgeo import gdal, ogr, osr

# Set input GeoTIFF file path
input_file = "./tmp/output/output_image1.tif"

# Set output KML file path
output_file = "./tmp/output/output1.kml"

# Open the GeoTIFF file
ds = gdal.Open(input_file)

# Get the GeoTransform
gt = ds.GetGeoTransform()

# Get the projection
proj = osr.SpatialReference()
proj.ImportFromWkt(ds.GetProjection())

# Create a KML datasource
driver = ogr.GetDriverByName("KML")
kml_ds = driver.CreateDataSource(output_file)

# Create a KML layer
kml_layer = kml_ds.CreateLayer("layer1", proj)

# Add a polygon to the KML layer
ring = ogr.Geometry(ogr.wkbLinearRing)
ring.AddPoint(gt[0], gt[3])
ring.AddPoint(gt[0] + gt[1] * ds.RasterXSize, gt[3])
ring.AddPoint(gt[0] + gt[1] * ds.RasterXSize, gt[3] + gt[5] * ds.RasterYSize)
ring.AddPoint(gt[0], gt[3] + gt[5] * ds.RasterYSize)
ring.AddPoint(gt[0], gt[3])
polygon = ogr.Geometry(ogr.wkbPolygon)
polygon.AddGeometry(ring)
feature = ogr.Feature(kml_layer.GetLayerDefn())
feature.SetGeometry(polygon)
kml_layer.CreateFeature(feature)

# Close the datasources
kml_ds = None
ds = None
