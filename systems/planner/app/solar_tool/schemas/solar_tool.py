from pydantic import BaseModel


class Site(BaseModel):
    latitude: float
    longitude: float
    power_budget: float | None = 130
    reliability_target: int | None = 98 # in percentile 95, 98 or 99
    
    class Config:
        orm_mode = True


class SolarToolRequestSchema(BaseModel):
    site: Site

    class Config:
        orm_mode = True


class SolarToolResponseSchema(BaseModel):
    number_of_solar_modules: int
    solar_pv_to_install_watts: float            # in watts
    #pv_module_cost: float
    number_of_batteries: int                    # The battery has to at least be the size of 1 day's energy consumption
    batteries_capacity_to_install_kWh: float    # in kWh
    max_output_angle: str    

    class Config:
        orm_mode = True
