import { format } from "date-fns";
import { Alert_Type } from "../generated";
import { TObject } from "../types";

const getTitleFromPath = (path: string) => {
    switch (path) {
        case "/":
            return "Home";
        case "/settings":
            return "Settings";
        case "/notification":
            return "Notification";
        case "/nodes":
            return "Nodes";
        case "/users":
            return "Users";
        case "/billing":
            return "Billing";
        case "/store":
            return "Module Store";
        default:
            return "Home";
    }
};

const getColorByType = (type: Alert_Type) => {
    switch (type) {
        case Alert_Type.Error:
            return "error";
        case Alert_Type.Warning:
            return "warning";
        default:
            return "success";
    }
};

const parseObjectInNameValue = (obj: any) => {
    let updatedObj: TObject[] = [];
    if (obj) {
        updatedObj = Object.keys(obj).map(key => {
            return {
                name: key,
                value:
                    key === "timestamp"
                        ? format(obj[key], "MMM dd HH:mm:ss")
                        : obj[key],
            };
        });

        let removeIndex = updatedObj
            .map(item => item?.name)
            .indexOf("__typename");
        ~removeIndex && updatedObj.splice(removeIndex, 1);
        removeIndex = updatedObj.map(item => item?.name).indexOf("id");
        ~removeIndex && updatedObj.splice(removeIndex, 1);
    }

    return updatedObj;
};

const uniqueObjectsArray = (name: string, list: TObject[]): TObject[] | [] => {
    const last =
        list.length > 0
            ? list.filter((item: TObject) => item.name !== name)
            : [];
    return last;
};

export {
    getColorByType,
    getTitleFromPath,
    uniqueObjectsArray,
    parseObjectInNameValue,
};
