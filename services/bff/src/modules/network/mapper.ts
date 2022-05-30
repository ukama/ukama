import { NetworkDto } from "./types";
import { INetworkMapper } from "./interface";
import { NETWORK_STATUS } from "../../constants";
import { MetricLatestValueRes } from "../../common/types";

class NetworkMapper implements INetworkMapper {
    dtoToDto = (res: MetricLatestValueRes): NetworkDto => {
        let uptime = 0;
        let status = NETWORK_STATUS.UNDEFINED;
        if (res) {
            uptime = parseFloat(res.value[1]);
            if (uptime > 0) {
                status = NETWORK_STATUS.ONLINE;
            } else {
                status = NETWORK_STATUS.DOWN;
            }
        }

        return { uptime: uptime, status: status };
    };
}
export default <INetworkMapper>new NetworkMapper();
