import { NetworkDto } from "./types";
import { MetricLatestValueRes, ParsedCookie } from "../../common/types";

export interface INetworkService {
    getNetworkStatus(cookie: ParsedCookie): Promise<NetworkDto>;
}

export interface INetworkMapper {
    dtoToDto(res: MetricLatestValueRes): NetworkDto;
}
