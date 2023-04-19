import os
from unittest.mock import patch, MagicMock
from app.elevation.services.sites_elevation import SitesElevation, Site, ElevationResponseSchema


class TestSitesElevation:
    @patch.object(SitesElevation, 'get_elevation_from_lon_lat')
    @patch.object(SitesElevation, 'get_raster_path')
    def test_get_elevations_success(self, mock_get_raster_path, mock_get_elevation_from_lon_lat):
        mock_get_raster_path.return_value = os.path.join(os.path.dirname(__file__), 'test_data', 'N00E010.hgt')
        mock_get_elevation_from_lon_lat.return_value = 100
        site1 = {"longitude": 10.0, "latitude": 0.0}
        site2 = {"longitude": 10.1, "latitude": 0.1}
        sites = [site1, site2]
        se = SitesElevation()
        elevations = se.get_elevations(sites)
        assert len(elevations) == 2
        assert elevations[0].longitude == site1['longitude']
        assert elevations[0].latitude == site1['latitude']
        assert elevations[0].elevation == 100
        assert elevations[1].longitude == site2['longitude']
        assert elevations[1].latitude == site2['latitude']
        assert elevations[1].elevation == 100

    @patch.object(SitesElevation, 'get_raster_path')
    def test_get_elevations_invalid_site(self, mock_get_raster_path):
        site1 = {"longitude": -200.0, "latitude": 0.0} # invalid longitude
        site2 = {"longitude": 10.1, "latitude": 200.0} # invalid latitude
        sites = [site1, site2]
        se = SitesElevation()
        elevations = se.get_elevations(sites)
        assert len(elevations) == 2
        assert elevations[0].elevation == 0
        assert elevations[1].elevation == 0
    
    @patch.object(SitesElevation, 'get_raster_path')
    def test_get_elevation_from_lon_lat_success(self, mock_get_raster_path):
        mock_ds = MagicMock()
        mock_ds.GetGeoTransform.return_value = (10.0, 0.000833333, 0.0, 0.0, -0.000833333, 10.0)
        mock_band = MagicMock()
        mock_band.ReadRaster.return_value = (0)
        mock_ds.GetRasterBand.return_value = mock_band
        mock_get_raster_path.return_value = os.path.join(os.path.dirname(__file__), 'test_data', 'N00E010.hgt')
        sites_elevation = SitesElevation()
        elevation = sites_elevation.get_elevation_from_lon_lat(10.0, 0.0)
        assert elevation == 0

    @patch.object(SitesElevation, 'get_raster_path')
    def test_get_elevation_from_lon_lat_invalid_lon_lat(self, mock_get_raster_path):
        sites_elevation = SitesElevation()
        elevation = sites_elevation.get_elevation_from_lon_lat(-200.0, 0.0) # invalid longitude
        assert elevation == 0
        elevation = sites_elevation.get_elevation_from_lon_lat(10.0, 200.0) # invalid latitude
        assert elevation == 0

    @patch("os.getenv", return_value="/path/to/hgt/files")
    @patch("core.config.get_config", return_value=MagicMock(HGT_FILES_PATH="/path/to/hgt/files"))
    def test_get_raster_path(self, mock_getenv, mock_get_config):
        sites_elevation = SitesElevation()
        sites_elevation.HGT_EXTENSION = ".hgt"
        
        # Test positive longitude and latitude values
        path = sites_elevation.get_raster_path(10, 20)
        assert path == "/path/to/hgt/files/N20E010.hgt"
        
        # Test negative longitude and latitude values
        path = sites_elevation.get_raster_path(-10, -20)
        assert path == "/path/to/hgt/files/S21W011.hgt"
        
        # Test longitude and latitude values at limits
        path = sites_elevation.get_raster_path(180, 90)
        assert path == "/path/to/hgt/files/N90E180.hgt"