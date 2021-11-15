import { BASE_URL } from "../constants/index";

export const URL = BASE_URL;

export const SERVER = {
    GET_CONNECTED_USERS: `${URL}/user/get_conneted_users`,
    GET_DATA_USAGE: `${URL}/data/data_usage`,
    GET_DATA_BILL: `${URL}/data/data_bill`,
    GET_ALERTS: `${URL}/alert/get_alerts`,
    GET_NODES: `${URL}/node/get_nodes`,
    GET_RESIDENTS: `${URL}/resident/get_residents`,
};
