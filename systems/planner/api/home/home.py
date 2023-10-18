from fastapi import APIRouter, Response, Depends

home_router = APIRouter()


@home_router.get("/health")
async def home():
    return Response(status_code=200)
