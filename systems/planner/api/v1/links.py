from fastapi import APIRouter, HTTPException
from app.links.services import SitesLinks
from app.links.schemas import LinksResponseSchema, LinksRequestSchema

links_router = APIRouter()


@links_router.post(
    "/links",
    response_model=LinksResponseSchema,
)
async def predict_links(request: LinksRequestSchema):
    try:
        response = SitesLinks().get_links(**request.dict())
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    else:
        return response
