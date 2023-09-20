import pytest
from fastapi import FastAPI, status
from httpx import AsyncClient
from api.v1.elevation import elevation_router

app = FastAPI()
app.include_router(elevation_router)


@pytest.fixture(scope="module")
def client():
    yield AsyncClient(app=app, base_url="http://test")


@pytest.mark.asyncio
async def test_predict_elevation(client):
    sites = [
        {"latitude": 51.86616553598713, "longitude": -2.2019728014010433},
        {"latitude": 51.849, "longitude": -2.2299},
    ]
    request_body = {"sites": sites}
    response_body = [
        {"latitude": 51.86616553598713, "longitude": -2.2019728014010433, "elevation":  0.0},
        {"latitude": 51.849, "longitude": -2.2299, "elevation":  0.0},
    ]
    response = await client.post("/elevation", json=request_body)

    assert response.status_code == status.HTTP_200_OK
    assert response.json() == response_body


@pytest.mark.asyncio
async def test_predict_elevation_with_empty_request_body(client):
    request_body = {}
    response = await client.post("/elevation", json=request_body)

    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
    assert response.json() == {
        "detail": [
            {
                "loc": ["body", "sites"],
                "msg": "field required",
                "type": "value_error.missing",
            }
        ]
    }


@pytest.mark.asyncio
async def test_predict_elevation_with_invalid_request_body(client):
    request_body = {"sites": "invalid"}
    response = await client.post("/elevation", json=request_body)

    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
    assert response.json() == {
        "detail": [
            {
                "loc": ["body", "sites"],
                "msg": "value is not a valid list",
                "type": "type_error.list",
            }
        ]
    }


@pytest.mark.asyncio
async def test_predict_elevation_with_missing_latitude(client):
    sites = [{"longitude": -122.676483}]
    request_body = {"sites": sites}
    response = await client.post("/elevation", json=request_body)

    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
    assert response.json() == {
        "detail": [
            {
                "loc": ["body", "sites", 0, "latitude"],
                "msg": "field required",
                "type": "value_error.missing",
            }
        ]
    }


@pytest.mark.asyncio
async def test_predict_elevation_with_missing_longitude(client):
    sites = [{"latitude": 45.523064}]
    request_body = {"sites": sites}
    response = await client.post("/elevation", json=request_body)

    assert response.status_code == status.HTTP_422_UNPROCESSABLE_ENTITY
    assert response.json() == {
        "detail": [
            {
                "loc": ["body", "sites", 0, "longitude"],
                "msg": "field required",
                "type": "value_error.missing",
            }
        ]
    }
