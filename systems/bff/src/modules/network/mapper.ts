import { MetricLatestValueRes } from "../../common/types";
import { NETWORK_STATUS } from "../../constants";
import { INetworkMapper } from "./interface";
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

class NetworkMapper implements INetworkMapper {
    dtoToSitesDto(res: SitesAPIResDto): SitesResDto {
        const sites: SiteDto[] = [];
        res.sites.forEach(site => {
            sites.push({
                id: site.id,
                name: site.name,
                networkId: site.network_id,
                isDeactivated: site.is_deactivated,
                createdAt: site.created_at,
            });
        });
        return {
            networkId: res.network_id,
            sites: sites,
        };
    }
    dtoToSiteDto(res: SiteAPIResDto): SiteDto {
        return {
            id: res.site.id,
            name: res.site.name,
            networkId: res.site.network_id,
            isDeactivated: res.site.is_deactivated,
            createdAt: res.site.created_at,
        };
    }
    dtoToNetworksDto(res: NetworksAPIResDto): NetworksResDto {
        const networks: NetworkDto[] = [];
        res.networks.forEach(network => {
            networks.push({
                id: network.id,
                name: network.name,
                orgId: network.org_id,
                isDeactivated: network.is_deactivated,
                createdAt: network.created_at,
            });
        });
        return {
            orgId: res.org_id,
            networks: networks,
        };
    }
    dtoToNetworkDto(res: NetworkAPIResDto): NetworkDto {
        return {
            id: res.network.id,
            name: res.network.name,
            orgId: res.network.org_id,
            isDeactivated: res.network.is_deactivated,
            createdAt: res.network.created_at,
        };
    }
    dtoToDto = (
        totalNodes: number,
        liveNodes: MetricLatestValueRes
    ): NetworkStatusDto => {
        let _liveNodes = 0;
        let status = NETWORK_STATUS.UNDEFINED;
        if (liveNodes) {
            _liveNodes = parseFloat(liveNodes.value[1]);
            if (_liveNodes > 0) {
                status = NETWORK_STATUS.ONLINE;
            } else {
                status = NETWORK_STATUS.DOWN;
            }
        }

        return { liveNode: _liveNodes, totalNodes: totalNodes, status: status };
    };
}
export default <INetworkMapper>new NetworkMapper();
