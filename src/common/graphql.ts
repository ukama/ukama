import "reflect-metadata";

export const GET_CONNECTED_USERS_QUERY = `
    query getUsers($data:TIME_FILTER!) {
        getConnectedUsers(filter:$data) {
            totalUser
            residentUsers
            guestUsers
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
                id
              title
              description
              totalUser
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
                id
              name
              usage
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
