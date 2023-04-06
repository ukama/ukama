from fastapi import APIRouter

from api.v1.coverage import coverage_router
from api.v1.elevation import elevation_router
from api.v1.links import links_router

router = APIRouter()
router.include_router(coverage_router, prefix="/api/v1", tags=["Coverage"])
router.include_router(elevation_router, prefix="/api/v1", tags=["Elevation"])
router.include_router(links_router, prefix="/api/v1", tags=["Links"])

__all__ = ["router"]
