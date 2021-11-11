import ApiMethods from "../../api";

export const getMethod = async (path: string, params: any, headers: any) => {
    const res = await ApiMethods.getData(path, params, headers);
    return res;
};
