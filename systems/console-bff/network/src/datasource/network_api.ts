import { RESTDataSource } from "@apollo/datasource-rest";
import { THeaders } from "../../../common/types";
import setupLogger from "../../../config/setupLogger";
import { SERVER } from "../../../constants/endpoints";
import { getHeaders } from "../../../utils";
import NetworkMapper from "./mapper";
import {
    AddNetworkInputDto,
    AddSiteInputDto,
    NetworkDto,
    NetworksResDto,
    NetworkStatusDto,
    SiteDto,
    SitesResDto,
} from "../types";
const logger = setupLogger("service");


export class NetworkApi extends RESTDataSource {
    getNetworkStatus = async (headers: THeaders): Promise<NetworkStatusDto> => {
        const resLiveNodes = await this.get(`${SERVER.ORG}/${headers.orgId}/metrics/live-nodes`,{headers: getHeaders(headers)});
        const resTotalNodes = await this.get(`${SERVER.ORG}/${headers.orgId}/nodes`,{headers: getHeaders(headers)});

        return NetworkMapper.dtoToDto(
            resTotalNodes.nodes.length,
            resLiveNodes.data.result[0]
        );
    };

    getNetworks = async (headers: THeaders): Promise<NetworksResDto> => {
        return this.get(`${SERVER.REGISTRY_NETWORKS_API_URL}`,{
            headers: getHeaders(headers),
            params: {
                org: headers.orgId,
            },}).then(res => 
            NetworkMapper.dtoToNetworksDto(res));
    };

    getNetwork = async (
        networkId: string,
        headers: THeaders
    ): Promise<NetworkDto> => {
        return this.get(`${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}`,{headers: getHeaders(headers)}).then(res => 
            NetworkMapper.dtoToNetworkDto(res));
    };

    getSites = async (
        networkId: string,
        headers: THeaders
    ): Promise<SitesResDto> => {
        return this.get(`${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites`,{headers: getHeaders(headers)}).then(res => 
            NetworkMapper.dtoToSitesDto(res));
        
    };

    getSite = async (
        siteId: string,
        networkId: string,
        headers: THeaders
    ): Promise<SiteDto> => {
        return this.get(`${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites/${siteId}`,{headers: getHeaders(headers)}).then(res => 
            NetworkMapper.dtoToSiteDto(res));
    };

    addNetwork = async (
        req: AddNetworkInputDto,
        headers: THeaders
    ): Promise<NetworkDto> => {
        return this.post(`${SERVER.REGISTRY_NETWORKS_API_URL}`,{
            headers: getHeaders(headers),
            body: req,
        }).then(res => 
            NetworkMapper.dtoToNetworkDto(res));
    };

    addSite = async (
        networkId: string,
        req: AddSiteInputDto,
        headers: THeaders
    ): Promise<SiteDto> => {
        return this.post(`${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites`, {
            headers: getHeaders(headers),
            body: req,
          }).then(res =>  NetworkMapper.dtoToSiteDto(res));
    };
}
