import "reflect-metadata";
import { gCall, beforeEachDeleteCall } from "../../../common/utils";
import { DELETE_USER_MUTATION } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: { id: "1342567fgcvh", success: true },
};
const id = "1342567fgcvh";

describe("Delete User", () => {
    beforeEachDeleteCall(
        "/user/delete_user?id=1342567fgcvh",
        nockResponse,
        200
    );
    it("delete users", async () => {
        const response = await gCall({
            source: DELETE_USER_MUTATION,
            variableValues: {
                input: id,
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
