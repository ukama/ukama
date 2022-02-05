import "reflect-metadata";

export const GET_CONNECTED_USERS_QUERY = `
    query getUsers($data:TIME_FILTER!) {
        getConnectedUsers(filter:$data) {
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

export const GET_RESIDENTS_QUERY = `
    query getResidents($input:PaginationDto!) {
        getResidents(data:$input) {
            residents{
                residents{
                    id
   		 	        name
                    status
				    eSimNumber
                    iccid
                    email
                    phone
                    dataPlan
 		   	        dataUsage
                    roaming
                }
                activeResidents
                totalResidents   
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

export const GET_USER_QUERY = `
    query getUsers($input:GetUserPaginationDto!) {
        getUsers(data:$input) {
            users{
                id
                name
                status
				eSimNumber
                iccid
                email
                phone
                dataPlan
                dataUsage
                roaming
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
    query getUser($input:String!) {
        getUser(id:$input) {
            id
            name
            status
            eSimNumber
            iccid
            email
            phone
            dataPlan
            dataUsage
            roaming
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

export const GET_NODE_PHYSICAL_HEALTH_QUERY = `
    query getNodePhysicalHealth {
        getNodePhysicalHealth { 
            temperature
            Memory
            cpu
            io
        }
    }
`;

export const GET_NODE_META_DATA_QUERY = `
    query getNodeMetaData {
        getNodeMetaData { 
            throughput
            usersAttached
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

export const GET_USERS_ATTACHED_METRICS_QUERY = `
query getUsersAttachedMetrics($filter: GRAPH_FILTER!) {
    getUsersAttachedMetrics(filter: $filter) {
     id
      users
      timestamp
    }
  }  
`;

export const GET_CPU_USAGE_METRICS_QUERY = `
query getCpuUsageMetrics($filter: GRAPH_FILTER!) {
    getCpuUsageMetrics(filter: $filter) {
      id
      usage
      timestamp
    }
  }  
`;

export const GET_NODE_RF_KPI_QUERY = `
query getNodeRFKPI($filter: GRAPH_FILTER!) {
    getNodeRFKPI(filter: $filter) {
      qam
      rfOutput
      rssi
      timestamp
    }
  }
`;

export const GET_TEMPERATURE_METRICS_QUERY = `
query getTemperatureMetrics($filter: GRAPH_FILTER!) {
    getTemperatureMetrics(filter: $filter) {
      id
      temperature
      timestamp
    }
  }
  
`;

export const GET_IO_METRICS_QUERY = `
query getIOMetrics($filter: GRAPH_FILTER!) {
    getIOMetrics(filter: $filter) {
      id
      input
      output
      timestamp
    }
  }  
`;

export const GET_THROUGHPUT_METRICS_QUERY = `
query getThroughputMetrics($filter: GRAPH_FILTER!) {
    getThroughputMetrics(filter: $filter) {
      uv
      pv
      amt
      time
    }
  }  
`;

export const GET_MEMORY_USAGE_METRICS_QUERY = `
    query getMemoryUsageMetrics($filter: GRAPH_FILTER!) {
        getMemoryUsageMetrics(filter: $filter) {
            id
            usage
            timestamp
        }
    }
  }
`;
