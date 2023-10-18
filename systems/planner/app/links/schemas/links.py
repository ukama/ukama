from typing import List
from pydantic import BaseModel


class Site(BaseModel):
    latitude: float
    longitude: float
    height: float | None = 10
    
    class Config:
        orm_mode = True


class LinksRequestSchema(BaseModel):
    sites: List[Site]

    class Config:
        orm_mode = True


class LinksResponseSchema(BaseModel):
    links: List[str]
    sites: List[Site]

    class Config:
        orm_mode = True
