import "reflect-metadata";
import { gCall, beforeEachPostCall } from "../../../common/utils";
import { DEACTIVATE_USER_MUTATION } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: { id: "1342567fgcvh", success: true },
};
const reqBody = {
    id: "1342567fgcvh",
};

describe("Deactivate User", () => {
    beforeEachPostCall("/user/deactivate_user", reqBody, nockResponse, 200);
    it("deactivate users", async () => {
        const response = await gCall({
            source: DEACTIVATE_USER_MUTATION,
            variableValues: {
                input: reqBody.id,
            },
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                deactivateUser: nockResponse.data,
            },
        });
    });
});
