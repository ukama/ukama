import os
import json
import requests
import numpy as np
import subprocess
from math import cos, radians, ceil

from app.solar_tool.schemas.solar_tool import Site, SolarToolResponseSchema
from app.solar_tool.enums.solar_tool import SolarToolEnum
from core import config

class SolarTool:
    def __init__(self):
        self.SOLAR_DATA_DIR = config.get_config().SOLAR_DATA_DIR

    def predict_solar_tools_requirements(self, site: Site) -> SolarToolResponseSchema:
        try:
            subprocess.check_call("mkdir -p " + self.SOLAR_DATA_DIR + "/")
            longitude, latitude, power_budget, reliability_target = site["longitude"], site["latitude"], site["power_budget"], site["reliability_target"]/100
            solar_data = self.get_solar_radiation(longitude, latitude) 

            sorted_three_days_avg_solar_data_values = sorted(list(solar_data['properties']['parameter']['ALLSKY_SFC_SW_DWN_MOVING_AVG'].values()))
            sorted_solar_data_values = sorted(list(solar_data['properties']['parameter']['ALLSKY_SFC_SW_DWN'].values()))

            total_daily_energy_consumption = 24*power_budget # Energy required by networking equipment for one day of uptime in Wh
            marginal_cost_battery_cycle =  round(SolarToolEnum.BATTERY_MODULE_COST.value/SolarToolEnum.BATTERY_MODULE_SIZE.value/SolarToolEnum.DEPTH_OF_DISCHARGE_PERCENTAGE.value/SolarToolEnum.BATTERY_CYCLE_LIFE.value, 2)  # Marginal cost of battery cycle ($/kWh)
            marginal_cost_battery_sun_hours = round(((SolarToolEnum.SOLAR_PANEL_COST_USD.value/SolarToolEnum.SOLAR_PANEL_SIZE_W.value)*1000)/(marginal_cost_battery_cycle*SolarToolEnum.SOLAR_PANEL_LIFETIME_YEARS.value*365), 2) # Above this insolation, it is cheaper to use solar; below, it is cheaper to use battery
            
            median_insolation = np.percentile(sorted_solar_data_values, 50)
            min_isolation_rel_target = round(np.percentile(sorted_three_days_avg_solar_data_values, (1-reliability_target)*100), 2)

            min_isolation = min([marginal_cost_battery_sun_hours, median_insolation])

            if min_isolation_rel_target > marginal_cost_battery_sun_hours:
                min_isolation = min_isolation_rel_target

            min_pv_size = total_daily_energy_consumption / SolarToolEnum.SOLAR_SYSTEM_TOTAL_DERATING_PERC.value / min_isolation

            total_solar_modules = ceil(min_pv_size / SolarToolEnum.SOLAR_PANEL_SIZE_W.value)
            solar_pv_to_install = total_solar_modules * SolarToolEnum.SOLAR_PANEL_SIZE_W.value # in watts

            pv_module_cost = total_solar_modules * SolarToolEnum.SOLAR_PANEL_COST_USD.value
            min_nominal_battery_required_kWh = round((3*total_daily_energy_consumption/SolarToolEnum.DEPTH_OF_DISCHARGE_PERCENTAGE.value - 2*solar_pv_to_install*min_isolation_rel_target*SolarToolEnum.SOLAR_SYSTEM_TOTAL_DERATING_PERC.value) / 1000, 2) 

            number_of_batteries = max([ceil(min_nominal_battery_required_kWh/SolarToolEnum.BATTERY_MODULE_SIZE.value), ceil(total_daily_energy_consumption/SolarToolEnum.DEPTH_OF_DISCHARGE_PERCENTAGE.value/1000/SolarToolEnum.BATTERY_MODULE_SIZE.value)])
            battery_capacity_to_install_kWh = round(number_of_batteries * SolarToolEnum.BATTERY_MODULE_SIZE.value, 2)

            battery_module_cost = number_of_batteries * SolarToolEnum.BATTERY_MODULE_COST.value
            estimated_capex = battery_module_cost+pv_module_cost+SolarToolEnum.BALANCE_SYSTEM_COST_USD.value

            return SolarToolResponseSchema(number_of_solar_modules=total_solar_modules, solar_pv_to_install_watts=solar_pv_to_install, number_of_batteries=number_of_batteries, batteries_capacity_to_install_kWh=battery_capacity_to_install_kWh)
        except Exception as ex:
            raise ex

    def get_3day_moving_avg(self, data):
        # Get the solar radiation values for the desired parameter and date range
        values = data["properties"]["parameter"]["ALLSKY_SFC_SW_DWN"]

        # Calculate the 3-day moving average for each date
        moving_averages = {}
        for i, date in enumerate(values.keys()):
            if i < 1:
                moving_averages[date] = round(values[date], 2)  # Use the same value for i < 1
            elif i == 1:
                prev_dates = list(values.keys())[i-1:i+1]
                moving_averages[date] = round((sum([values[d] for d in prev_dates]) / len(prev_dates)), 2)  # Average of first and second values
            else:
                prev_dates = list(values.keys())[i-2:i+1]
                moving_averages[date] = round((sum([values[d] for d in prev_dates]) / len(prev_dates)), 2)

        return moving_averages

    def get_solar_radiation_from_api(self, longitude, latitude):
        #now = datetime.now()
        #end_date = now.strftime("%Y%m%d")
        start_date = "19840101"    # hardcoding it for now
        end_date = "20221231"       
        
        url = f"https://power.larc.nasa.gov/api/temporal/daily/point?parameters=ALLSKY_SFC_SW_DWN&community=RE&longitude={longitude}&latitude={latitude}&start={start_date}&end={end_date}&format=JSON"
        response = requests.get(url)
        if response.status_code == 200:
            data = response.json()
            return data
        else:
            raise Exception(f"Error: {response.status_code} - {response.reason}")
    
    def calculate_ranges(self, longitude, latitude):
        area = 1878.4 # Area covered by the NASA API from one point of longitude and latitude in meters
        # Calculate the range of longitude in degrees
        d_lon = (2 * area) / (111320 * cos(radians(latitude)))
        
        # Calculate the range of latitude in degrees
        d_lat = (2 * area) / 111320
        
        # Calculate the longitude and latitude ranges
        longitude_range = (round(longitude - d_lon/2, 3), round(longitude + d_lon/2, 3))
        latitude_range = (round(latitude - d_lat/2, 3), round(latitude + d_lat/2, 3))
        
        return longitude_range, latitude_range
    
    def search_saved_files(self, longitude_range, latitude_range):
        # Search for saved files in the longitude and latitude range
        for filename in os.listdir(self.SOLAR_DATA_DIR):
            if filename.endswith(".json"):
                lon_lat_str = filename[:-5] # remove .json extension
                lon_str, lat_str = lon_lat_str.split("_")
                lon = float(lon_str)
                lat = float(lat_str)
                if (lon >= longitude_range[0] and lon <= longitude_range[1] and
                    lat >= latitude_range[0] and lat <= latitude_range[1]):
                    return filename
        
        return None
    
    def save_json_into_file(self, output_data, output_filename):
        # Dump the output to the file
        out_file = open(self.SOLAR_DATA_DIR + "/" + output_filename, 'w')
        try:
            json.dump(output_data, out_file)
        finally:
            out_file.close()

    def get_solar_radiation(self, longitude, latitude):
        # Calculate the longitude and latitude ranges
        longitude_range, latitude_range = self.calculate_ranges(longitude, latitude)
        saved_solar_data_filename = self.search_saved_files(longitude_range, latitude_range)
        solar_data = {}
        if saved_solar_data_filename == None:  # Means the file is not found in the range of given longitude and latitude, so fetch from the api and save it.
            # The longitude and latitude is within the range of the 1x1 block, so save the data in the file using range
            output_filename = f"{longitude_range[0]}_{latitude_range[0]}.json"

            # call the API to fetch the data
            output = self.get_solar_radiation_from_api(longitude, latitude)
            output["properties"]["parameter"]["ALLSKY_SFC_SW_DWN"] = {k: v for k, v in output["properties"]["parameter"]["ALLSKY_SFC_SW_DWN"].items() if v != -999}
            
            # Calculate the moving averages and add them to the output
            moving_averages = self.get_3day_moving_avg(output)
            output["properties"]["parameter"]["ALLSKY_SFC_SW_DWN_MOVING_AVG"] = moving_averages
            solar_data = output
            self.save_json_into_file(output, output_filename)
        else:    
            # Read from the file if it longitude and latitude is in range
            with open(self.SOLAR_DATA_DIR + "/" + saved_solar_data_filename, "r") as f:
                solar_data = json.load(f)
        return solar_data
