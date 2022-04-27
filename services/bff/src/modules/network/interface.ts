import { NETWORK_TYPE } from "../../constants";
import { NetworkDto, NetworkResponse } from "./types";

export interface INetworkService {
    getNetwork(filter: NETWORK_TYPE): Promise<NetworkDto>;
}

export interface INetworkMapper {
    dtoToDto(res: NetworkResponse): NetworkDto;
}
