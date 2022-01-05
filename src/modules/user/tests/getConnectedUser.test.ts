import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_CONNECTED_USERS_QUERY } from "../../../common/graphql";
import { HEADER, TIME_FILTER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        totalUser: 8,
    },
};

describe("Get Connected Users", () => {
    beforeEachGetCall("/user/get_conneted_users?0=WEEK", nockResponse, 200);
    it("get connected users", async () => {
        const response = await gCall({
            source: GET_CONNECTED_USERS_QUERY,
            variableValues: {
                data: TIME_FILTER.WEEK,
            },
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                getConnectedUsers: nockResponse.data,
            },
        });
    });
});
