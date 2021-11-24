import { ErrorType } from "../common/types";
import { HttpStatusCode } from "./codes";
import { HTTP400Error } from "./http400.error";
import { HTTP401Error } from "./http401.error";
import { HTTP404Error } from "./http404.error";
import { HTTP500Error } from "./http500.error";
import Messages from "./messages";

export {
    HTTP400Error,
    HTTP500Error,
    HTTP401Error,
    HttpStatusCode,
    HTTP404Error,
    Messages,
};

export const axiosErrorHandler = (error: any): ErrorType => {
    let res: ErrorType;
    if (error.response) {
        // The request was made and the server responded with a status code
        res = {
            code: error.response.status,
            message: error.response.statusText,
        };
    } else if (error.request) {
        // The request was made but no response was received
        res = {
            code: error.request.status,
            message: error.request.statusText,
        };
    } else {
        // Something happened in setting up the request that triggered an Error
        res = {
            code: 400,
            message: error.message,
        };
    }
    return res;
};

export const checkError = (error: any): boolean => {
    if (error.code && error.message) return true;
    return false;
};
