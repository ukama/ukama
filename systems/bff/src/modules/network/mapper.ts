import { NetworkDto } from "./types";
import { INetworkMapper } from "./interface";
import { NETWORK_STATUS } from "../../constants";
import { MetricLatestValueRes } from "../../common/types";

class NetworkMapper implements INetworkMapper {
    dtoToDto = (
        totalNodes: number,
        liveNodes: MetricLatestValueRes,
    ): NetworkDto => {
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
