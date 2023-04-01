from fastapi import APIRouter, HTTPException
from app.coverage.services import SitesCoverage
from app.coverage.schemas import CoverageResponseSchema, CoverageRequestSchema

coverage_router = APIRouter()


@coverage_router.post(
    "/coverage",
    response_model=CoverageResponseSchema,
)
async def predict_coverage(request: CoverageRequestSchema):
    try:
        response = SitesCoverage().calculate_coverage(**request.dict())
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    else:
        return response
