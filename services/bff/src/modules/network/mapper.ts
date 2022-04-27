import { INetworkMapper } from "./interface";
import { NetworkResponse, NetworkDto } from "./types";

class NetworkMapper implements INetworkMapper {
    dtoToDto = (res: NetworkResponse): NetworkDto => {
        return res.data;
    };
}
export default <INetworkMapper>new NetworkMapper();
