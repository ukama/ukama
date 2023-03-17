import { MetricLatestValueRes, ParsedCookie } from "../../common/types";
import {
    NetworkAPIResDto,
    NetworkDto,
    NetworksAPIResDto,
    NetworksResDto,
    NetworkStatusDto,
    SiteAPIResDto,
    SiteDto,
    SitesAPIResDto,
    SitesResDto,
} from "./types";

export interface INetworkService {
    getNetworkStatus(cookie: ParsedCookie): Promise<NetworkStatusDto>;
}

export interface INetworkMapper {
    dtoToDto(
        totalNodes: number,
        liveNodes: MetricLatestValueRes,
    ): NetworkStatusDto;
    dtoToNetworksDto(res: NetworksAPIResDto): NetworksResDto;
    dtoToNetworkDto(res: NetworkAPIResDto): NetworkDto;
    dtoToSitesDto(res: SitesAPIResDto): SitesResDto;
    dtoToSiteDto(res: SiteAPIResDto): SiteDto;
}
