from typing import List
from pydantic import BaseModel, Field
from sqlalchemy import Column, Float, Integer
from sqlalchemy.ext.declarative import declarative_base


Base = declarative_base()

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

class PopulationDataResponse(BaseModel):
    url: str
    population_covered: float
    total_boxes_covered: int
    class Config:
        orm_mode = True

class InterferenceDataResponse(BaseModel):
    url: str
    class Config:
        orm_mode = True

class CoverageResponseSchema(BaseModel):
    north: float
    east: float
    west: float
    south: float
    url: str
    population_data: dict[str, PopulationDataResponse]
    interference_data: dict[str, InterferenceDataResponse]
    class Config:
        orm_mode = True
