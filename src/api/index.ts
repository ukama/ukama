import axios from "axios";
import { GETDataDto } from "../common/types";
import setupLogger from "../config/logger";

const logger = setupLogger("Api Methods");

class ApiMethods {
    getData = async (req: GETDataDto) => {
        let res;

        try {
            res = await axios.get(req.path, {
                params: req.params,
                headers: {
                    ...req.headers,
                    Accept: "*/*",
                },
            });
        } catch (error) {
            res = null;
            logger.error(error);
        }
        return res;
    };
}

export default new ApiMethods();
