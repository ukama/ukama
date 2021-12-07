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
        return await axios({
            method: type,
            url: path,
            data: body,
            headers: headers,
            params: params,
        }).catch(err => {
            throw err;
        });
    };
}

export default new ApiMethods();
