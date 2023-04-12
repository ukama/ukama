from fastapi import APIRouter, HTTPException
from typing import List
from app.elevation.services import SitesElevation
from app.elevation.schemas import ElevationResponseSchema, ElevationRequestSchema

elevation_router = APIRouter()


@elevation_router.post(
    "/elevation",
    response_model=List[ElevationResponseSchema],
)
async def predict_coverage(request: ElevationRequestSchema):
    try:
        response = SitesElevation().get_elevations(**request.dict())
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    else:
        return response
