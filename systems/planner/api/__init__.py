from fastapi import APIRouter

from api.v1.coverage import coverage_router

router = APIRouter()
router.include_router(coverage_router, prefix="/api/v1", tags=["Coverage"])

__all__ = ["router"]
