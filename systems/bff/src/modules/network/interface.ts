import { MetricLatestValueRes } from "../../common/types";
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
// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface INetworkService {}

export interface INetworkMapper {
    dtoToDto(
        totalNodes: number,
        liveNodes: MetricLatestValueRes
    ): NetworkStatusDto;
    dtoToNetworksDto(res: NetworksAPIResDto): NetworksResDto;
    dtoToNetworkDto(res: NetworkAPIResDto): NetworkDto;
    dtoToSitesDto(res: SitesAPIResDto): SitesResDto;
    dtoToSiteDto(res: SiteAPIResDto): SiteDto;
}
