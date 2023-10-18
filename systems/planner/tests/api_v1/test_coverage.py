from unittest.mock import MagicMock

from fastapi.testclient import TestClient

from api.v1.coverage import coverage_router
from app.coverage.services import SitesCoverage


client = TestClient(coverage_router)


def test_predict_coverage_success():
    # Mock the SitesCoverage service's `calculate_coverage` method to return a response
    SitesCoverage.calculate_coverage = MagicMock(return_value={"north": 40.1, "east": -70.2, "west": -71.3, "south": 39.4, "url": "/c/output/test13.png"})

    # Define a request payload
    payload = {
        "mode": "simple",
        "sites": [
            {"latitude": 39.9, "longitude": -75.2},
            {"latitude": 40.0, "longitude": -75.0},
        ],
    }

    # Send a request to the endpoint
    response = client.post("/coverage", json=payload)

    # Assert that the response has a 200 status code
    assert response.status_code == 200

    # Assert that the response body matches the expected response
    assert response.json() == {"north": 40.1, "east": -70.2, "west": -71.3, "south": 39.4, "url": "/c/output/test13.png"}
