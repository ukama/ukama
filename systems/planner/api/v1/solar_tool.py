from fastapi import APIRouter, HTTPException
from app.solar_tool.services import SolarTool
from app.solar_tool.schemas import SolarToolResponseSchema, SolarToolRequestSchema

solar_tool_router = APIRouter()


@solar_tool_router.post(
    "/solar_tools",
    response_model=SolarToolResponseSchema,
)
async def predict_solar_tools(request: SolarToolRequestSchema):
    try:
        response = SolarTool().predict_solar_tools_requirements(**request.dict())
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    else:
        return response
