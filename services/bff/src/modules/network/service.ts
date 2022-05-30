import { Service } from "typedi";
import { NetworkDto } from "./types";
import { checkError } from "../../errors";
import { INetworkService } from "./interface";
import { catchAsyncIOMethod } from "../../common";
import { SERVER } from "../../constants/endpoints";
import { API_METHOD_TYPE } from "../../constants";
import { ParsedCookie } from "../../common/types";
import setupLogger from "../../config/setupLogger";
import NetworkMapper from "./mapper";
const logger = setupLogger("service");
@Service()
export class NetworkService implements INetworkService {
    getNetworkStatus = async (cookie: ParsedCookie): Promise<NetworkDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/metrics/live-status`,
            headers: cookie.header,
        });
        if (checkError(res)) {
            logger.error(res);
            throw new Error(res.message);
        }

        return NetworkMapper.dtoToDto(res.data.result[0]);
    };
}
