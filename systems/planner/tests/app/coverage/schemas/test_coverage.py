import pytest
from app.coverage.schemas import Site, CoverageRequestSchema, CoverageResponseSchema
from pydantic import ValidationError

class TestSite:
    def test_valid_site(self):
        site_dict = {'latitude': 40.7128, 'longitude': -74.0060}
        site = Site(**site_dict)
        assert site.latitude == site_dict['latitude']
        assert site.longitude == site_dict['longitude']
        assert site.transmitter_height == 25

    def test_site_with_height(self):
        site_dict = {'latitude': 40.7128, 'longitude': -74.0060, 'transmitter_height': 30}
        site = Site(**site_dict)
        assert site.transmitter_height == site_dict['transmitter_height']

    def test_site_with_invalid_latitude(self):
        site_dict = {'latitude': 'invalid', 'longitude': -74.0060}
        with pytest.raises(ValidationError):
            Site(**site_dict)

    def test_site_with_invalid_longitude(self):
        site_dict = {'latitude': 40.7128, 'longitude': 'invalid'}
        with pytest.raises(ValidationError):
            Site(**site_dict)

    def test_site_with_invalid_height(self):
        site_dict = {'latitude': 40.7128, 'longitude': -74.0060, 'transmitter_height': 'invalid'}
        with pytest.raises(ValidationError):
            Site(**site_dict)

    def test_site_with_missing_latitude(self):
        site_dict = {'longitude': -74.0060}
        with pytest.raises(ValidationError):
            Site(**site_dict)

    def test_site_with_missing_longitude(self):
        site_dict = {'latitude': 40.7128}
        with pytest.raises(ValidationError):
            Site(**site_dict)

    def test_site_with_missing_height(self):
        site_dict = {'latitude': 40.7128, 'longitude': -74.0060, 'transmitter_height': None}
        site = Site(**site_dict)
        assert site.transmitter_height is None

class TestCoverageSchemas:
    def test_valid_coverage_request_schema(self):
        # Define a valid coverage request
        valid_coverage_request = {
            "mode": "receive_power",
            "sites": [
                {
                    "latitude": 51.5072,
                    "longitude": -0.1276,
                    "transmitter_height": 30
                }
            ]
        }
        # Validate the coverage request schema
        coverage_request = CoverageRequestSchema(**valid_coverage_request)
        assert coverage_request.mode == "receive_power"
        assert len(coverage_request.sites) == 1
        assert isinstance(coverage_request.sites[0], Site)
        assert coverage_request.sites[0].latitude == 51.5072
        assert coverage_request.sites[0].longitude == -0.1276
        assert coverage_request.sites[0].transmitter_height == 30


    def test_invalid_coverage_request_schema(self):
        # Define an invalid coverage request with a missing mode
        invalid_coverage_request = {
            "sites": [
                {
                    "latitude": 51.5072,
                    "longitude": -0.1276,
                    "transmitter_height": 30
                }
            ]
        }
        # Validate the coverage request schema
        with pytest.raises(ValueError):
            CoverageRequestSchema(**invalid_coverage_request)

        # Define an invalid coverage request with a site missing latitude
        invalid_coverage_request = {
            "mode": "receive_power",
            "sites": [
                {
                    "longitude": -0.1276,
                    "transmitter_height": 30
                }
            ]
        }
        # Validate the coverage request schema
        with pytest.raises(ValueError):
            CoverageRequestSchema(**invalid_coverage_request)

        # Define an invalid coverage request with a site missing longitude
        invalid_coverage_request = {
            "mode": "receive_power",
            "sites": [
                {
                    "latitude": 51.5072,
                    "transmitter_height": 30
                }
            ]
        }
        # Validate the coverage request schema
        with pytest.raises(ValueError):
            CoverageRequestSchema(**invalid_coverage_request)

        # Define an invalid coverage request with a site containing an invalid latitude
        invalid_coverage_request = {
            "mode": "receive_power",
            "sites": [
                {
                    "latitude": "invalid_latitude",
                    "longitude": -0.1276,
                    "transmitter_height": 30
                }
            ]
        }
        # Validate the coverage request schema
        with pytest.raises(ValueError):
            CoverageRequestSchema(**invalid_coverage_request)

        # Define an invalid coverage request with a site containing an invalid longitude
        invalid_coverage_request = {
            "mode": "receive_power",
            "sites": [
                {
                    "latitude": 51.5072,
                    "longitude": "invalid_longitude",
                    "transmitter_height": 30
                }
            ]
        }
        # Validate the coverage request schema
        with pytest.raises(ValueError):
            CoverageRequestSchema(**invalid_coverage_request)


class TestCoverageResponseSchema:
    @classmethod
    def setup_class(self):
        """setup any state specific to the execution of the given class (which
        usually contains tests).
        """
        self.response_data = {
            "north": 51.5074,
            "east": -0.1278,
            "west": -0.1478,
            "south": 51.4874,
            "url": "https://example.com"
        }


    def test_coverage_response_schema(self):
        # Test valid response data
        response = CoverageResponseSchema(**self.response_data)
        assert response.north == self.response_data["north"]
        assert response.east == self.response_data["east"]
        assert response.west == self.response_data["west"]
        assert response.south == self.response_data["south"]
        assert response.url == self.response_data["url"]

        # Test missing required fields
        with pytest.raises(ValueError):
            CoverageResponseSchema()

        # Test invalid data types
        invalid_data = self.response_data.copy()
        invalid_data["north"] = "invalid"
        with pytest.raises(ValueError):
            CoverageResponseSchema(**invalid_data)
        