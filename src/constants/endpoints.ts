import env from "../config/env";

export const BASE_URL = env.BASE_URL;

export const SERVER = {
    GET_CONNECTED_USERS: `${BASE_URL}/user/get_conneted_users`,
    GET_DATA_USAGE: `${BASE_URL}/data/data_usage`,
};
