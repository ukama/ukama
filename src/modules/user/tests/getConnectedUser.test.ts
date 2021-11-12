import "reflect-metadata";
import { gCall, nockCall } from "../../../test/utils";
import { GET_CONNECTED_USERS_QUERY } from "../../../test/graphql";
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
    beforeEach(() => {
        nockCall("/user/get_conneted_users?0=WEEK", nockResponse);
    });
    it("get connected users", async () => {
        const response = await gCall({
            source: GET_CONNECTED_USERS_QUERY,
            variableValues: {
                data: TIME_FILTER.WEEK,
            },
        });
        expect(response).toMatchObject({
            data: {
                getConnectedUsers: nockResponse.data,
            },
        });
    });
});
