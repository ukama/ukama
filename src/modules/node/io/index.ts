import { AxiosResponse } from "axios";
import ApiMethods from "../../../api";
import { PaginationDto } from "../../../common/types";
import { API_METHOD_TYPE } from "../../../constants";
import { SERVER } from "../../../constants/endpoints";

class NodeIOMethods {
    getNodesMethod = async (
        params: PaginationDto
    ): Promise<AxiosResponse<any, any> | null> => {
        const res = await ApiMethods.fetch({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_NODES,
            params: params,
        });
        return res;
    };
}
export default new NodeIOMethods();
