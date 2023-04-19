import pytest
from app.links.schemas import Site, LinksResponseSchema, LinksRequestSchema
from pydantic import ValidationError

class TestSite:
    def test_create_site(self):
        site = Site(latitude=37.7749, longitude=-122.4194, height=12)
        assert site.latitude == 37.7749
        assert site.longitude == -122.4194
        assert site.height == 12


    def test_create_site_without_height(self):
        site = Site(latitude=37.7749, longitude=-122.4194)
        assert site.latitude == 37.7749
        assert site.longitude == -122.4194
        assert site.height == 10


    def test_create_site_with_invalid_latitude(self):
        with pytest.raises(ValidationError):
            Site(latitude="91a", longitude=-122.4194, height=12)


    def test_create_site_with_invalid_longitude(self):
        with pytest.raises(ValidationError):
            Site(latitude=37.7749, longitude="91a", height=12)


    def test_create_site_with_invalid_height(self):
        with pytest.raises(ValidationError):
            Site(latitude=37.7749, longitude=-122.4194, height="91a")

class TestLinksRequestSchema:
    def test_create_links_request_schema(self):
        sites = [Site(latitude=37.7749, longitude=-122.4194, height=12)]
        request = LinksRequestSchema(sites=sites)
        assert len(request.sites) == 1
        assert request.sites[0].latitude == 37.7749
        assert request.sites[0].longitude == -122.4194
        assert request.sites[0].height == 12


    def test_create_links_request_schema_with_multiple_sites(self):
        sites = [
            Site(latitude=37.7749, longitude=-122.4194, height=12),
            Site(latitude=40.7128, longitude=-74.0060, height=20)
        ]
        request = LinksRequestSchema(sites=sites)
        assert len(request.sites) == 2
        assert request.sites[0].latitude == 37.7749
        assert request.sites[0].longitude == -122.4194
        assert request.sites[0].height == 12
        assert request.sites[1].latitude == 40.7128
        assert request.sites[1].longitude == -74.0060
        assert request.sites[1].height == 20

    def test_create_links_request_schema_with_null_sites(self):
        with pytest.raises(ValidationError):
            LinksRequestSchema(sites=None)


class TestLinksResponseSchema:
    def test_create_links_response_schema(self):
        sites = [Site(latitude=37.7749, longitude=-122.4194, height=12)]
        response = LinksResponseSchema(links=[], sites=sites)
        assert len(response.sites) == 1
        assert response.sites[0].latitude == 37.7749
        assert response.sites[0].longitude == -122.4194
        assert response.sites[0].height == 12