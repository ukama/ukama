import pytest
from fastapi import FastAPI, status
from httpx import AsyncClient

from api.v1.links import links_router

app = FastAPI()
app.include_router(links_router)


@pytest.fixture(scope="module")
def client():
    yield AsyncClient(app=app, base_url="http://test")

@pytest.mark.asyncio
async def test_predict_links(client):
    sites = [
         {"latitude": 51.86616553598713, "longitude": -2.2019728014010433, "height": 30},
        {"latitude": 51.849, "longitude": 2.2299, "height": 60},
    ]

    request_body = {"sites": sites}
    response_body = {
        "links": [
            "(51.86616553598713, -2.2019728014010433) -> (51.849, 2.2299)"
        ],
        "sites": [
            {"height":  40.5, "latitude": 51.86616553598713, "longitude": -2.2019728014010433},
            {"height":  40.5, "latitude": 51.849, "longitude": 2.2299},
        ]
    }
    response = await client.post("/links", json=request_body)

    assert response.status_code == status.HTTP_200_OK
    assert response.json() == response_body
    
@pytest.mark.asyncio
async def test_predict_links_with_empty_request_body(client):
    request_body = {}
    response = await client.post("/links", json=request_body)

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
async def test_predict_links_with_invalid_request_body(client):
    request_body = {"sites": "invalid"}
    response = await client.post("/links", json=request_body)

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
async def test_predict_links_with_missing_latitude(client):
    sites = [{"longitude": -122.676483}]
    request_body = {"sites": sites}
    response = await client.post("/links", json=request_body)

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
async def test_predict_links_with_missing_longitude(client):
    sites = [{"latitude": 45.523064}]
    request_body = {"sites": sites}
    response = await client.post("/links", json=request_body)

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
