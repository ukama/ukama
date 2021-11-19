import "reflect-metadata";
import { gCall, beforeEachCall } from "../../../common/utils";
import { GET_CONNECTED_USERS_QUERY } from "../../../common/graphql";
import { TIME_FILTER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        totalUser: 8,
        residentUsers: 5,
        guestUsers: 3,
    },
};

describe("Get Connected Users", () => {
    beforeEachCall("/user/get_conneted_users?0=WEEK", nockResponse, 200);
    it("get connected users", async () => {
        const response = await gCall({
            source: GET_CONNECTED_USERS_QUERY,
            variableValues: {
                data: TIME_FILTER.WEEK,
            },
            contextValue: {
                req: {
                    headers: {
                        authorization: "test",
                    },
                },
            },
        });
        expect(response).toMatchObject({
            data: {
                getConnectedUsers: nockResponse.data,
            },
        });
    });
});
