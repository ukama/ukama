# build a schema using pydantic
from pydantic import BaseModel

class Site(BaseModel):
    lon: float
    lat: float
    elv: float
    
    class Config:
        orm_mode = True

class SiteCoverage(BaseModel):
    lon: float
    lat: float
    freq: int
    rt: int | None = -90
    pm: int | None = 1
    txh: int | None = 25
    erp: int | None = 0
    radius: int | None = 30
    res: int | None = 600
    sdf_filename: str | None = ""
    rel_itm_model: int | None = None
    rx_threshold: int | None = None
    dbm: bool | None = False 
    output_filename: str | None = "test"
    class Config:
        orm_mode = True

class SiteCoverageReponse(BaseModel):
    output_filename: str
    output_type: str
    png_path: str
    ppm_path: str
    dcf_path: str
    east_cordinate: float
    west_cordinate: float
    north_cordinate: float
    south_cordinate: float
    
    class Config:
        orm_mode = True

class ErrorReponse(BaseModel):
    error_message: str
    
    class Config:
        orm_mode = True