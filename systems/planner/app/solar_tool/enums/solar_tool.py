from . import BaseEnum


class SolarToolEnum(BaseEnum):
    DEPTH_OF_DISCHARGE_PERCENTAGE = 0.80    # In percentage, usable battery capacity
    BATTERY_CYCLE_LIFE = 6000             # Number of cycles at above depth of discharge
    BATTERY_MODULE_SIZE = 2.4             # Nominal battery size
    BATTERY_MODULE_COST = 1400            # In dollars             
    SOLAR_PANEL_SIZE_W = 400              # Watts
    SOLAR_PANEL_COST_USD = 300            # Include cost of racking, etc in US dollars
    SOLAR_PANEL_LIFETIME_YEARS = 20       # In years
    BALANCE_SYSTEM_COST_USD = 2000        # MPPT charger, PoE Switch, DC-DC converter, enclosure, breakers, computer, wiring, etc in US dollars
    SOLAR_SYSTEM_TOTAL_DERATING_PERC = 0.82 # Solar System Total Derating in percent, Product of derating factors for solar (ratio of energy stored in battery to energy from solar panels)