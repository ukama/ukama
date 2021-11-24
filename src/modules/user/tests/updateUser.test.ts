import "reflect-metadata";
import { gCall, beforeEachPostCall } from "../../../common/utils";
import { POST_UPDATE_USER_MUTATION } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = { status: "success", data: { success: true } };
const reqBody = {
    firstName: "ali",
};

describe("Post Update Users", () => {
    beforeEachPostCall("/user/update_user", reqBody, nockResponse, 200);
    it("post update users", async () => {
        const response = await gCall({
            source: POST_UPDATE_USER_MUTATION,
            variableValues: {
                input: reqBody,
            },
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                updateUser: nockResponse.data,
            },
        });
    });
});
