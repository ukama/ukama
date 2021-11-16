import axios from "axios";
import { URL } from "../constants/endpoints";
import { ApiMethodDataDto } from "../common/types";

class ApiMethods {
    constructor() {
        axios.create({
            timeout: 10000,
        });
    }
    fetch = async (req: ApiMethodDataDto) => {
        const { headers, path, params, type, body } = req;
        return await axios({
            method: type,
            url: `${URL}/${path}`,
            data: body,
            headers: headers,
            params: params,
        }).catch(() => {
            return null;
        });
    };
}

export default new ApiMethods();
