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

export const GET_Alerts_QUERY = `
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
