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
        const resLiveNodes = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/metrics/live-nodes`,
            headers: cookie.header,
        });
        const resTotalNodes = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: `${SERVER.ORG}/${cookie.orgId}/nodes`,
            headers: cookie.header,
        });
        if (checkError(resLiveNodes)) {
            logger.error(resLiveNodes);
            throw new Error(resLiveNodes.message);
        }
        if (checkError(resTotalNodes)) {
            logger.error(resTotalNodes);
            throw new Error(resTotalNodes.message);
        }

        return NetworkMapper.dtoToDto(
            resTotalNodes.nodes.length,
            resLiveNodes.data.result[0]
        );
    };
}
