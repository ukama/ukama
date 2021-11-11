import "reflect-metadata";
import { gCall } from "../../../test/utils";
import { GET_CONNECTED_USERS_QUERY } from "../../../test/graphql";
import nock from "nock";
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
        nock("http://127.0.0.1:3000")
            .get("/user/get_conneted_users?0=WEEK")
            .reply(200, nockResponse);
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
