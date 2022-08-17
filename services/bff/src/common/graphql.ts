import "reflect-metadata";

export const GET_CONNECTED_USERS_QUERY = `
    query getUsers($orgId: String!, $data: TIME_FILTER!) {
        getConnectedUsers(filter: $data, orgId: $orgId) {
            totalUser
        }
    }
`;

export const GET_DATA_USAGE_QUERY = `
    query getDataUsage($data:TIME_FILTER!) {
        getDataUsage(filter: $data) {
            id
            dataConsumed
            dataPackage
        }
    }
`;

export const GET_DATA_BILL_QUERY = `
    query getDataBill($data:DATA_BILL_FILTER!) {
        getDataBill(filter: $data) {
            id
            dataBill
            billDue
        }
    }
`;

export const GET_ALERTS_QUERY = `
    query getAlerts($input:PaginationDto!) {
        getAlerts(data:$input) {
            alerts{
                id
                type,
                title,
                description,
                alertDate
            }
            meta{
                page
                count
                pages
                size
            }
        }
    }
`;

export const GET_NODES_QUERY = `
    query getNodes($input:PaginationDto!) {
        getNodes(data:$input) {
            nodes{
                nodes{
                    id
                    title
                    description
                    status
                    totalUser
              }
                activeNodes
                totalNodes 
              }
              meta{
                page
                count
                pages
                size
              }
        }
    }
`;

export const POST_ACTIVATE_USER_MUTATION = `
    mutation activateUser($input:ActivateUserDto!) {
        activateUser(data:$input) {
            success
        }
    }
`;

export const GET_ESIM_QUERY = `
    query getEsims {
        getEsims{
            esim
            active
        }
    }
`;

export const POST_ADD_NODE_MUTATION = `
    mutation addNode($input:AddNodeDto!) {
        addNode(data:$input) {
            success
        }
    }
`;

export const GET_CURRENT_BILL = `
    query getCurrentBill {
        getCurrentBill {
            bill {
             id
             name
             dataUsed
              rate
             subtotal
            }
          total
          dueDate
          billMonth
        }
    }
`;

export const GET_BILL_HISTORY = `
    query getBillHistory {
        getBillHistory {
            id
            description
            date
            totalUsage
            subtotal
            invoice
        }
    }
`;

export const GET_NETWORK_QUERY = `
    query getNetwork($data:NETWORK_TYPE!) {
        getNetwork(filter: $data) {
            id
            status
            description
        }
    }
`;

export const POST_UPDATE_USER_MUTATION = `
    mutation updateUser($input:UpdateUserDto!) {
        updateUser(data:$input) {
            id
            name
            phone
            sim
            email
        }
    }
`;

export const DEACTIVATE_USER_MUTATION = `
    mutation deactivateUser($input:String!) {
        deactivateUser(id:$input){
            id
            success
        }
    }
`;

export const POST_UPDATE_NODE_MUTATION = `
    mutation updateNode($input:UpdateNodeDto!) {
        updateNode(data:$input) {
            id
            name
            serialNo
        }
    }
`;

export const DELETE_NODE_MUTATION = `
    mutation deleteNode($input:String!) {
        deleteNode(id:$input){
            id
            success
        }
    }
`;

export const GET_USER_BY_ID_QUERY = `
    query getUser($userInput: UserInput!) {
        getUser(userInput: $userInput) {
            id
            status
            name
            eSimNumber
            iccid
            email
            phone
            roaming
            dataPlan
            dataUsage
        }
    }
`;

export const GET_NODE_DETAILS_QUERY = `
    query getNodeDetails {
        getNodeDetails { 
            id
            modelType
            serial
            macAddress
            osVersion
            manufacturing
            ukamaOS
            hardware
            description
        }
    }
`;

export const GET_NODE_NETWORK = `
    query getNodeNetwork {
        getNodeNetwork { 
            id
            status
            description
        }
    }
`;
