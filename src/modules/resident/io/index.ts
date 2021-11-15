import { AxiosResponse } from "axios";
import ApiMethods from "../../../api";
import { PaginationDto } from "../../../common/types";
import { SERVER } from "../../../constants/endpoints";

export const getResidentsMethod = async (
    params: PaginationDto
): Promise<AxiosResponse<any, any> | null> => {
    const res = await ApiMethods.getData({
        path: SERVER.GET_RESIDENTS,
        params: params,
    });
    return res;
};
