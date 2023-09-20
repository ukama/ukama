import pytest
from unittest.mock import MagicMock
from app.links.services.sites_links import SitesLinks

class TestSitesLinks:

    def test_generate_links(self):
        # Define test inputs
        locations = [(10.0, 20.0), (30.0, 40.0), (50.0, 60.0)]
        # Define expected output
        expected_links = [((10.0, 20.0), (30.0, 40.0)), ((30.0, 40.0), (50.0, 60.0))]
        # Instantiate SitesLinks class and call the generate_links method
        sites_links = SitesLinks()
        links = sites_links.generate_links(locations)
        # Check that the output matches the expected output
        assert links == expected_links

    def test_predict_heights_from_links(self, monkeypatch):
        # Define test inputs
        links = [((10.0, 20.0), (30.0, 40.0)), ((30.0, 40.0), (50.0, 60.0))]
        # Define expected output
        expected_towers_with_heights = {
            (30.0, 40.0): (204.5, 10.0),
            (10.0, 20.0): (138.0, 10.0),
            (50.0, 60.0): (86.5, 10.0)
        }
        # Create a mock object for the SitesElevation class
        mock_sites_elevation = MagicMock()
        mock_sites_elevation.get_elevation_from_lon_lat.return_value = 10.0
        # Set the mock object to be returned by the SitesElevation instance
        monkeypatch.setattr('app.links.services.sites_links.SitesElevation', lambda: mock_sites_elevation)
        # Instantiate SitesLinks class and call the predict_heights_from_links method
        sites_links = SitesLinks()
        towers_with_heights = sites_links.predict_heights_from_links(links)
        # Check that the output matches the expected output
        assert towers_with_heights == expected_towers_with_heights

    @pytest.fixture
    def mock_sites_elevation(monkeypatch):
        mock_sites_elevation = MagicMock()
        mock_sites_elevation.get_elevation_from_lon_lat.return_value = 10.0
        monkeypatch.setattr('app.links.services.sites_links.SitesElevation', lambda: mock_sites_elevation)
        yield mock_sites_elevation

    def test_get_link_status_returns_true(mock_sites_elevation):
        sites_links = SitesLinks()
        towerA_height = (30.0, 20.0)
        towerB_height = (40.0, 30.0)
        tower1_loc = (51.86616553598713, -2.2019728014010433)
        tower2_loc = (51.849, -2.2299)
        status = sites_links.get_link_status(towerA_height, towerB_height, tower1_loc, tower2_loc)
        assert status

    def test_get_link_status_returns_false(mock_sites_elevation):
        sites_links = SitesLinks()
        towerA_height = (30.0, 20.0)
        towerB_height = (10.0, 5.0)
        tower1_loc = (10.0, 20.0)
        tower2_loc = (30.0, 40.0)
        assert sites_links.get_link_status(towerA_height, towerB_height, tower1_loc, tower2_loc) == False

    def test_get_link_status_returns_false_if_tower1_height_is_greater_than_tower2_height(mock_sites_elevation):
        sites_links = SitesLinks()
        towerA_height = (50.0, 30.0)
        towerB_height = (40.0, 20.0)
        tower1_loc = (10.0, 20.0)
        tower2_loc = (30.0, 40.0)
        status = sites_links.get_link_status(towerA_height, towerB_height, tower1_loc, tower2_loc)
        assert status == False
