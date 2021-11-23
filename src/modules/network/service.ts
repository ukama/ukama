import { Service } from "typedi";
import { NetworkDto, NetworkResponse } from "./types";
import { INetworkService } from "./interface";
import NetworkMapper from "./mapper";
import { catchAsyncIOMethod } from "../../common";
import { SERVER } from "../../constants/endpoints";
import { API_METHOD_TYPE, NETWORK_TYPE } from "../../constants";

@Service()
export class NetworkService implements INetworkService {
    getNetwork = async (filter: NETWORK_TYPE): Promise<NetworkDto> => {
        const res = await catchAsyncIOMethod<NetworkResponse>({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NETWORK,
            params: `${filter}`,
        });
        const network = NetworkMapper.dtoToDto(res);
        return network;
    };
}
