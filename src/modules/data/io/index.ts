import { AxiosResponse } from "axios";
import ApiMethods from "../../../api";
import {
    API_METHOD_TYPE,
    DATA_BILL_FILTER,
    TIME_FILTER,
} from "../../../constants";
import { SERVER } from "../../../constants/endpoints";

class DataIOMethods {
    getDataUsageMethod = async (
        filter: TIME_FILTER
    ): Promise<AxiosResponse<any, any> | null> => {
        const res = await ApiMethods.fetch({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_DATA_USAGE,
            params: `${filter}`,
        });
        return res;
    };

    getDataBillMethod = async (
        filter: DATA_BILL_FILTER
    ): Promise<AxiosResponse<any, any> | null> => {
        const res = await ApiMethods.fetch({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_DATA_BILL,
            params: `${filter}`,
        });
        return res;
    };
}

export default new DataIOMethods();
