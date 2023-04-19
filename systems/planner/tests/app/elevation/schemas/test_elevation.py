import pytest
from app.elevation.schemas import Site, ElevationResponseSchema, ElevationRequestSchema
from pydantic import ValidationError

class TestSite:
    def test_create_site(self):
        site = Site(latitude=52.5200, longitude=13.4050)
        assert site.latitude == 52.5200
        assert site.longitude == 13.4050

    def test_create_site_with_invalid_latitude(self):
        site_dict = {'latitude': 'invalid', 'longitude': -74.0060}
        with pytest.raises(ValidationError):
            Site(**site_dict)


    def test_create_site_with_invalid_longitude(self):
        site_dict = {'latitude': -74.0060, 'longitude': 'invalid'}
        with pytest.raises(ValidationError):
            Site(**site_dict)

class TestCoverageSchemas:
    def test_elevation_request_schema_with_valid_data(self):
        data = {
            "sites": [
                {"latitude": 42.3601, "longitude": -71.0589},
                {"latitude": 40.7128, "longitude": -74.0060}
            ]
        }
        schema = ElevationRequestSchema(**data)
        assert schema.sites[0].latitude == 42.3601
        assert schema.sites[0].longitude == -71.0589
        assert schema.sites[1].latitude == 40.7128
        assert schema.sites[1].longitude == -74.0060


    def test_elevation_request_schema_with_invalid_data(self):
        data = {
            "sites": [
                {"latitude": 42.3601, "longitude": "invalid"},
                {"latitude": 40.7128, "longitude": -74.0060}
            ]
        }
        with pytest.raises(ValueError):
            ElevationRequestSchema(**data)


    def test_elevation_request_schema_with_missing_data(self):
        data = {}
        with pytest.raises(ValueError):
            ElevationRequestSchema(**data)


class TestElevationResponseSchema:
    @classmethod
    def setup_class(self):
        """setup any state specific to the execution of the given class (which
        usually contains tests).
        """
        self.elevation_response_data = {
            "latitude": 51.5074,
            "longitude": 0.1278,
            "elevation": 35.0
        }

    def test_elevation_response_schema(self):
        response = ElevationResponseSchema.parse_obj(self.elevation_response_data)
        assert response.latitude == self.elevation_response_data["latitude"]
        assert response.longitude == self.elevation_response_data["longitude"]
        assert response.elevation == self.elevation_response_data["elevation"]


    def test_elevation_response_schema_missing_fields(self):
        del self.elevation_response_data["latitude"]
        with pytest.raises(ValueError):
            ElevationResponseSchema.parse_obj(self.elevation_response_data)


    def test_elevation_response_schema_invalid_data_types(self):
        self.elevation_response_data["latitude"] = "invalid"
        with pytest.raises(ValueError):
            ElevationResponseSchema.parse_obj(self.elevation_response_data)
