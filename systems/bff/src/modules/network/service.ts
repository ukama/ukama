import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { THeaders } from "../../common/types";
import setupLogger from "../../config/setupLogger";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError } from "../../errors";
import { getHeaders } from "../../utils";
import { INetworkService } from "./interface";
import NetworkMapper from "./mapper";
import {
    AddNetworkInputDto,
    AddSiteInputDto,
    NetworkDto,
    NetworksResDto,
    NetworkStatusDto,
    SiteDto,
    SitesResDto,
} from "./types";
const logger = setupLogger("service");
@Service()
export class NetworkService implements INetworkService {
    getNetworkStatus = async (headers: THeaders): Promise<NetworkStatusDto> => {
        const resLiveNodes = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${headers.orgId}/metrics/live-nodes`,
            headers: getHeaders(headers),
        });
        const resTotalNodes = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${headers.orgId}/nodes`,
            headers: getHeaders(headers),
        });
        if (checkError(resLiveNodes)) {
            logger.error(resLiveNodes);
            throw new Error(resLiveNodes.message);
        }
        if (checkError(resTotalNodes)) {
            logger.error(resTotalNodes);
            throw new Error(resTotalNodes.message);
        }

        return NetworkMapper.dtoToDto(
            resTotalNodes.nodes.length,
            resLiveNodes.data.result[0]
        );
    };

    getNetworks = async (headers: THeaders): Promise<NetworksResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_NETWORKS_API_URL}`,
            params: {
                org: headers.orgId,
            },
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        return NetworkMapper.dtoToNetworksDto(res);
    };

    getNetwork = async (
        networkId: string,
        headers: THeaders
    ): Promise<NetworkDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}`,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        return NetworkMapper.dtoToNetworkDto(res);
    };

    getSites = async (
        networkId: string,
        headers: THeaders
    ): Promise<SitesResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites`,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        return NetworkMapper.dtoToSitesDto(res);
    };

    getSite = async (
        siteId: string,
        networkId: string,
        headers: THeaders
    ): Promise<SiteDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites/${siteId}`,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        return NetworkMapper.dtoToSiteDto(res);
    };

    addNetwork = async (
        req: AddNetworkInputDto,
        headers: THeaders
    ): Promise<NetworkDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.REGISTRY_NETWORKS_API_URL}`,
            body: req,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        return NetworkMapper.dtoToNetworkDto(res);
    };

    addSite = async (
        networkId: string,
        req: AddSiteInputDto,
        headers: THeaders
    ): Promise<SiteDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: `${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites`,
            body: req,
            headers: getHeaders(headers),
        });

        if (checkError(res)) throw new Error(res.message);
        return NetworkMapper.dtoToSiteDto(res);
    };
}
