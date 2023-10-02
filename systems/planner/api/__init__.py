from fastapi import APIRouter
from api import *
from api.v1.coverage import coverage_router
from api.v1.elevation import elevation_router
from api.v1.links import links_router
from api.v1.solar_tool import solar_tool_router

router = APIRouter()
router.include_router(coverage_router, prefix="/api/v1", tags=["Coverage"])
router.include_router(elevation_router, prefix="/api/v1", tags=["Elevation"])
router.include_router(links_router, prefix="/api/v1", tags=["Links"])
router.include_router(solar_tool_router, prefix="/api/v1", tags=["SolarTool"])

__all__ = ["router"]
