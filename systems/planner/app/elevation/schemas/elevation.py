from typing import List
from pydantic import BaseModel


class Site(BaseModel):
    latitude: float
    longitude: float
    
    class Config:
        orm_mode = True


class ElevationRequestSchema(BaseModel):
    sites: List[Site]

    class Config:
        orm_mode = True


class ElevationResponseSchema(BaseModel):
    latitude: float
    longitude: float
    elevation: float

    class Config:
        orm_mode = True