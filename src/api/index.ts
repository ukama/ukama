import axios from "axios";
import setupLogger from "../config/logger";

const logger = setupLogger("Api Methods");

class ApiMethods {
    getData = async (path: string, params?: any, headers?: any) => {
        const res = await axios.get(path, {
            params,
            headers: {
                ...headers,
                Accept: "*/*",
            },
        });
        return res;
    };
    putData = async (path: string, params?: any, headers?: any, body?: any) => {
        const res = await axios.put(path, body, {
            params,
            headers: {
                ...headers,
                Accept: "*/*",
            },
        });
        return res;
    };
    postData = async (
        path: string,
        params?: any,
        headers?: any,
        body?: any
    ) => {
        let res;
        try {
            res = await axios.post(path, body, {
                params,
                headers: {
                    ...headers,
                    Accept: "*/*",
                },
            });
        } catch (error: any) {
            res = error.response.data.message;
            logger.error(error);
        }
        return res;
    };
    patchData = async (
        path: string,
        params?: any,
        headers?: any,
        body?: any
    ) => {
        const res = await axios.patch(path, body, {
            params,
            headers: {
                ...headers,
                Accept: "*/*",
            },
        });
        return res;
    };
    deleteData = async (path: string, params?: any, headers?: any) => {
        const res = await axios.delete(path, {
            params,
            headers: {
                ...headers,
                Accept: "*/*",
            },
        });
        return res;
    };
}

export default new ApiMethods();
