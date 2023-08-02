from typing import List
from pydantic import BaseModel, Field


class Site(BaseModel):
    latitude: float
    longitude: float
    transmitter_height: int | None = 25
    coverage_radius: float | None = 30
    
    class Config:
        orm_mode = True


class CoverageRequestSchema(BaseModel):
    mode: str
    population_coverage: bool | None = True
    provide_interference_data: bool | None = True
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
    population_data: dict[str, PopulationDataResponse] | None
    interference_data: dict[str, InterferenceDataResponse] | None
    class Config:
        orm_mode = True
