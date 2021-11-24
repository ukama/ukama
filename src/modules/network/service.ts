import { Service } from "typedi";
import { NetworkDto } from "./types";
import { INetworkService } from "./interface";
import NetworkMapper from "./mapper";
import { catchAsyncIOMethod } from "../../common";
import { SERVER } from "../../constants/endpoints";
import { API_METHOD_TYPE, NETWORK_TYPE } from "../../constants";
import { checkError } from "../../errors";

@Service()
export class NetworkService implements INetworkService {
    getNetwork = async (filter: NETWORK_TYPE): Promise<NetworkDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NETWORK,
            params: `${filter}`,
        });
        if (checkError(res)) throw new Error(res.message);

        const network = NetworkMapper.dtoToDto(res);
        return network;
    };
}
