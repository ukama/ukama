import { AxiosResponse } from "axios";
import ApiMethods from "../../../api";
import { DATA_BILL_FILTER, TIME_FILTER } from "../../../constants";
import { SERVER } from "../../../constants/endpoints";

export const getDataUsageMethod = async (
    filter: TIME_FILTER
): Promise<AxiosResponse<any, any> | null> => {
    const res = await ApiMethods.getData({
        path: SERVER.GET_DATA_USAGE,
        params: `${filter}`,
    });
    return res;
};

export const getDataBillMethod = async (
    filter: DATA_BILL_FILTER
): Promise<AxiosResponse<any, any> | null> => {
    const res = await ApiMethods.getData({
        path: SERVER.GET_DATA_BILL,
        params: `${filter}`,
    });
    return res;
};
