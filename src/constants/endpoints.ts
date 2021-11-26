import { BASE_URL } from "../constants/index";

export const URL = BASE_URL;

export const SERVER = {
    GET_CONNECTED_USERS: `user/get_conneted_users`,
    GET_DATA_USAGE: `data/data_usage`,
    GET_DATA_BILL: `data/data_bill`,
    GET_ALERTS: `alert/get_alerts`,
    GET_NODES: `node/get_nodes`,
    GET_RESIDENTS: `resident/get_residents`,
    GET_ESIMS: `esims/get_esims`,
    POST_ACTIVE_USER: `user/active_user`,
    GET_USERS: `user/get_users`,
    POST_ADD_NODE: `node/add_node`,
    GET_CURRENT_BILL: `bill/get_current_bill`,
    GET_BILL_HISTORY: `bill/get_bill_history`,
    GET_NETWORK: `network/get_network`,
    POST_UPDATE_USER: `user/update_user`,
    DELETE_USER: `user/delete_user`,
};
