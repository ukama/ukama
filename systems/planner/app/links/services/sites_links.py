import math
import networkx as nx
from haversine import haversine, Unit
from typing import List

from app.links.schemas.links import Site, LinksResponseSchema
from app.elevation.services.sites_elevation import SitesElevation

class SitesLinks:
    def __init__(self):
        self.sites_elevation = SitesElevation()
        self.fast_solution = True

    def get_links(self, sites: List[Site]) -> LinksResponseSchema:
        try:
            locations = [(site['latitude'], site['longitude']) for site in sites]
            links = self.generate_links(locations)
            links_response = [f"{link[0]} -> {link[1]}" for link in links]
            towerHeights = self.predict_heights_from_links(links)
            new_sites = []
            for (lat, lon), (total_height, loc_height) in towerHeights.items():
                site = Site(latitude=lat, longitude=lon, height=(total_height-loc_height))
                new_sites.append(site)
            return LinksResponseSchema(links=links_response, sites=new_sites)
        except Exception as e:
            raise e
    

    def generate_links(self, locations):
        graph = nx.Graph()
        for i, site1 in enumerate(locations):
            for j, site2 in enumerate(locations):
                if i < j:
                    # calculate distance between two locations 
                    distance = haversine(site1, site2, unit=Unit.METERS)
                    #print(f"{i}, {j} -> {distance}")
                    graph.add_edge(i, j, weight=distance)
        mst = nx.minimum_spanning_tree(graph)
        linksArr = []
        for edge in mst.edges(data=True):
            #print(locations[edge[0]], locations[edge[1]], edge[2]['weight'])
            linksArr.append((locations[edge[0]], locations[edge[1]]))
        return linksArr
        

    def predict_heights_from_links(self, links):
        default_height_in_meters = 10 #meters
        towers_with_heights = {}
        for i, (tower1_loc, tower2_loc) in enumerate(links): # filling up default heights
            loc_elevation1 = self.sites_elevation.get_elevation_from_lon_lat(tower1_loc[1], tower1_loc[0])
            towers_with_heights[tower1_loc] = (loc_elevation1 + default_height_in_meters, loc_elevation1)
            loc_elevation2 = self.sites_elevation.get_elevation_from_lon_lat(tower2_loc[1], tower2_loc[0])
            towers_with_heights[tower2_loc] = (loc_elevation2 + default_height_in_meters, loc_elevation2)
        for i, (tower1_loc, tower2_loc) in enumerate(links):
            while not self.get_link_status(towers_with_heights[tower1_loc], towers_with_heights[tower2_loc], tower1_loc, tower2_loc):
                total_height1, loc_elevation1 = towers_with_heights[tower1_loc]
                total_height2, loc_elevation2 = towers_with_heights[tower2_loc]
                total_height1 += 0.5
                total_height2 += 0.5
                towers_with_heights[tower1_loc] = (total_height1, loc_elevation1)
                towers_with_heights[tower2_loc] = (total_height2, loc_elevation2)
        return towers_with_heights

    def get_link_status(self, towerA_height, towerB_height, tower1_loc, tower2_loc) -> bool:
            distance_A_B_km = haversine(tower1_loc, tower2_loc)
            # Fresnel zone clearance
            freq= 5 #Frequency in GHz

            xFresnelClearance = 60
            for i in range(1, 101):
                fraction = i / 101
                xlat = tower1_loc[0] + fraction * (tower2_loc[0] - tower1_loc[0])
                xlon = tower1_loc[1] + fraction * (tower2_loc[1] - tower1_loc[1])
                location_height = self.sites_elevation.get_elevation_from_lon_lat(xlon, xlat)
                distance_X_A = haversine((xlat, xlon), tower1_loc)
                xLoS = ((i * (towerB_height[0] - towerA_height[0])) / 100) + towerA_height[0]
                xFresnelFactor = 17.3 * math.sqrt((distance_X_A * (distance_A_B_km - distance_X_A)) / (distance_A_B_km * freq))
                xFresnelFactor = xFresnelFactor * xFresnelClearance / 100
                xFresnelZone = xLoS - xFresnelFactor
                xObstm = 0 if (xFresnelZone - location_height) >= 0 else xFresnelFactor + location_height
                if xObstm > 0:
                    return False
            #print(f"There is a line of sight between the two towers. {tower1_loc}, {tower2_loc}")
            return True
