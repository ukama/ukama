from fastapi import APIRouter

from api.v1.coverage import coverage_router
from api.v1.elevation import elevation_router

router = APIRouter()
router.include_router(coverage_router, prefix="/api/v1", tags=["Coverage"])
router.include_router(elevation_router, prefix="/api/v1", tags=["Elevation"])

__all__ = ["router"]
