import axios from "axios";
import { ApiMethodDataDto } from "../common/types";

class ApiMethods {
    constructor() {
        axios.create({
            timeout: 10000,
        });
    }
    fetch = async (req: ApiMethodDataDto) => {
        const { headers, path, params, type, body } = req;
        return axios({
            method: type,
            url: path,
            data: body,
            headers: headers,
            params: params,
            timeout: 10000,
        }).catch(err => {
            throw err;
        });
    };
}

export default new ApiMethods();
