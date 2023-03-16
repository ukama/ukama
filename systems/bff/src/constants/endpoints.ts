import { BASE_URL, DEV_URL } from "../constants/index";

export const SERVER = {
    REGISTRY_NODE_API_URL: `${DEV_URL}/v1/nodes`,
    GET_CONNECTED_USERS: `${BASE_URL}/user/get_conneted_users`,
    GET_DATA_USAGE: `${BASE_URL}/data/data_usage`,
    GET_DATA_BILL: `${BASE_URL}/data/data_bill`,
    GET_ALERTS: `${BASE_URL}/alert/get_alerts`,
    GET_NODES: `${BASE_URL}/node/get_nodes`,
    GET_RESIDENTS: `${BASE_URL}/resident/get_residents`,
    GET_ESIMS: `${BASE_URL}/esims/get_esims`,
    POST_ACTIVE_USER: `${BASE_URL}/user/active_user`,
    GET_USERS: `${BASE_URL}/user/get_users`,
    POST_ADD_NODE: `${BASE_URL}/node/add_node`,
    GET_CURRENT_BILL: `${BASE_URL}/bill/get_current_bill`,
    GET_BILL_HISTORY: `${BASE_URL}/bill/get_bill_history`,
    GET_NETWORK: `${BASE_URL}/network/get_network`,
    POST_UPDATE_USER: `${BASE_URL}/user/update_user`,
    POST_DEACTIVATE_USER: `${BASE_URL}/user/deactivate_user`,
    POST_UPDATE_NODE: `${BASE_URL}/node/update_node`,
    POST_DELETE_NODE: `${BASE_URL}/node/delete_node`,
    GET_USER: `${BASE_URL}/user/get_user`,
    ORG: `${DEV_URL}/orgs`,
    GET_SOFTWARE_LOGS: `${BASE_URL}/software_logs`,
    GET_NODE_APPS: `${BASE_URL}/node_apps`,
    GET_NODE_DETAIL: `${BASE_URL}/node/node_details`,
    GET_NODE_META_DATA: `${BASE_URL}/node/meta_data`,
    GET_NODE_PHYSICAL_HEALTH: `${BASE_URL}/node/physical_health`,
    GET_NODE_RF_KPI: `${BASE_URL}/node/rf_kpis`,
    GET_NODE_NETWORK: `${BASE_URL}/node/get_network`,
    GET_THROUGHPUT_METRICS: `${BASE_URL}/metrics/throughput`,
    GET_USERS_ATTACHED_METRICS: `${BASE_URL}/metrics/user`,
    GET_CPU_USAGE_METRICS: `${BASE_URL}/metrics/cpu`,
    GET_TEMPERATURE_METRICS: `${BASE_URL}/metrics/temperature`,
    GET_IO_METRICS: `${BASE_URL}/metrics/io`,
    GET_MEMORY_USAGE_METRICS: `${BASE_URL}/metrics/memory`,
    GET_IDENTITY: `https://kratos-admin.dev.ukama.com/admin/identities`,
};

export const getMetricUri = (
    orgId: string,
    nodeId: string,
    endpoint: string
): string => `${SERVER.ORG}/${orgId}/nodes/${nodeId}/metrics/${endpoint}`;
