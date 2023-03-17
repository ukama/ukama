import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { ParsedCookie } from "../../common/types";
import setupLogger from "../../config/setupLogger";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError } from "../../errors";
import { INetworkService } from "./interface";
import NetworkMapper from "./mapper";
import {
    NetworkDto,
    NetworksResDto,
    NetworkStatusDto,
    SiteDto,
    SitesResDto,
} from "./types";
const logger = setupLogger("service");
@Service()
export class NetworkService implements INetworkService {
    getNetworkStatus = async (
        cookie: ParsedCookie,
    ): Promise<NetworkStatusDto> => {
        const resLiveNodes = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/metrics/live-nodes`,
            headers: cookie.header,
        });
        const resTotalNodes = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/nodes`,
            headers: cookie.header,
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
            resLiveNodes.data.result[0],
        );
    };

    getNetworks = async (cookie: ParsedCookie): Promise<NetworksResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_NETWORKS_API_URL}`,
            params: {
                org: cookie.orgId,
            },
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        return NetworkMapper.dtoToNetworksDto(res);
    };

    getNetwork = async (
        networkId: string,
        cookie: ParsedCookie,
    ): Promise<NetworkDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        return NetworkMapper.dtoToNetworkDto(res);
    };

    getSites = async (
        networkId: string,
        cookie: ParsedCookie,
    ): Promise<SitesResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        return NetworkMapper.dtoToSitesDto(res);
    };

    getSite = async (
        siteId: string,
        networkId: string,
        cookie: ParsedCookie,
    ): Promise<SiteDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.REGISTRY_NETWORKS_API_URL}/${networkId}/sites/${siteId}`,
            headers: cookie.header,
        });

        if (checkError(res)) throw new Error(res.message);
        return NetworkMapper.dtoToSiteDto(res);
    };
}
