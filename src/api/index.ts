import axios from "axios";
import setupLogger from "../config/logger";

const logger = setupLogger("Api Methods");

class ApiMethods {
    getData = async (path: string, params?: any, headers?: any) => {
        let res;

        try {
            res = await axios.get(path, {
                params,
                headers: {
                    ...headers,
                    Accept: "*/*",
                },
            });
        } catch (error: any) {
            res = null;
            logger.error(error);
        }
        return res;
    };
}

export default new ApiMethods();
