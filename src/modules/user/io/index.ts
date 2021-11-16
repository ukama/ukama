import { AxiosResponse } from "axios";
import ApiMethods from "../../../api";
import { API_METHOD_TYPE, TIME_FILTER } from "../../../constants";
import { SERVER } from "../../../constants/endpoints";

class UserIOMethods {
    getUsersMethod = async (
        filter: TIME_FILTER
    ): Promise<AxiosResponse<any, any> | null> => {
        const res = await ApiMethods.fetch({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_CONNECTED_USERS,
            params: `${filter}`,
        });
        return res;
    };
}

export default new UserIOMethods();
