import ApiMethods from "../../../api";

export const getDataUsageMethod = async (
    path: string,
    params: any,
    headers: any
) => {
    const res = await ApiMethods.getData(path, params, headers);
    return res;
};
