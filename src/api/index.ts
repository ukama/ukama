// import axios from "axios";
// import { ApiMethodDataDto } from "../common/types";
// import setupLogger from "../config/logger";

// const logger = setupLogger("Api Methods");

// class ApiMethods {
//     getData = async (req: ApiMethodDataDto) => {
//         let res;

//         try {
//             res = await axios.get(req.path, {
//                 params: req.params,
//                 headers: {
//                     ...req.headers,
//                     Accept: "*/*",
//                 },
//             });
//         } catch (error) {
//             res = null;
//             logger.error(error);
//         }
//         return res;
//     };
// }

import axios from "axios";
import setupLogger from "../config/logger";
import { URL } from "../constants/endpoints";
import { ApiMethodDataDto } from "../common/types";
const logger = setupLogger("Api Methods");
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
        }).catch(error => {
            logger.error(error);
            return null;
        });
    };
}

export default new ApiMethods();
