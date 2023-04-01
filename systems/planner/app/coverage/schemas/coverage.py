from typing import List
from pydantic import BaseModel


class Site(BaseModel):
    latitude: float
    longitude: float
    transmitter_height: int | None = 25
    
    class Config:
        orm_mode = True


class CoverageRequestSchema(BaseModel):
    mode: str
    sites: List[Site]

    class Config:
        orm_mode = True


class CoverageResponseSchema(BaseModel):
    north: float
    east: float
    west: float
    south: float
    url: str

    class Config:
        orm_mode = True