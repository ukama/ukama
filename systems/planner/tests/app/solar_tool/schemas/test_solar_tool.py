import pytest
from app.solar_tool.schemas import Site, SolarToolResponseSchema, SolarToolRequestSchema
from pydantic import ValidationError

class TestSite:
    def test_create_site(self):
        site = Site(latitude=37.7749, longitude=-122.4194, power_budget=120, reliability_target=98)
        assert site.latitude == 37.7749
        assert site.longitude == -122.4194
        assert site.power_budget == 120
        assert site.reliability_target == 98


    def test_create_site_without_power_budget(self):
        site = Site(latitude=37.7749, longitude=-122.4194, reliability_target=98)
        assert site.latitude == 37.7749
        assert site.longitude == -122.4194
        assert site.power_budget == 130
        assert site.reliability_target == 98

    def test_create_site_without_reliability_target(self):
        site = Site(latitude=37.7749, longitude=-122.4194)
        assert site.latitude == 37.7749
        assert site.longitude == -122.4194
        assert site.power_budget == 130
        assert site.reliability_target == 98

    def test_create_site_with_invalid_latitude(self):
        with pytest.raises(ValidationError):
            Site(latitude="91a", longitude=-122.4194, height=12)


    def test_create_site_with_invalid_longitude(self):
        with pytest.raises(ValidationError):
            Site(latitude=37.7749, longitude="91a", height=12)


class TestSolarToolRequestSchema:
    def test_create_solar_tool_request_schema(self):
        site = Site(latitude=37.7749, longitude=-122.4194, power_budget=120, reliability_target=98)
        request = SolarToolRequestSchema(site=site)
        assert request.site.latitude == 37.7749
        assert request.site.longitude == -122.4194
        assert request.site.power_budget == 120
        assert request.site.reliability_target == 98


    def test_create_solar_tool_request_schema_with_null_sites(self):
        with pytest.raises(ValidationError):
            SolarToolRequestSchema(site=None)


class TestSolarToolResponseSchema:
    def test_create_solar_tool_response_schema(self):
        response = SolarToolResponseSchema(number_of_solar_modules=5, solar_pv_to_install_watts=35.2, number_of_batteries=23, batteries_capacity_to_install_kWh=212.2, max_output_angle= "10 degrees north")
        assert response.number_of_solar_modules == 5
        assert response.solar_pv_to_install_watts == 35.2
        assert response.number_of_batteries == 23
        assert response.batteries_capacity_to_install_kWh == 212.2
        assert response.max_output_angle == "10 degrees north"
