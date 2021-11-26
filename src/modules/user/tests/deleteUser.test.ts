import "reflect-metadata";
import { gCall, beforeEachPostCall } from "../../../common/utils";
import { DELETE_USER_MUTATION } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: { id: "1342567fgcvh", success: true },
};
const reqBody = {
    id: "1342567fgcvh",
};

describe("Delete User", () => {
    beforeEachPostCall("/user/delete_user", reqBody, nockResponse, 200);
    it("delete users", async () => {
        const response = await gCall({
            source: DELETE_USER_MUTATION,
            variableValues: {
                input: reqBody.id,
            },
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                deleteUser: nockResponse.data,
            },
        });
    });
});
