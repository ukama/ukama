import "reflect-metadata";
import { gCall, beforeEachPostCall } from "../../../common/utils";
import { POST_ACTIVATE_USER_MUTATION } from "../../../common/graphql";

const nockResponse = { status: "success", data: { success: true } };
const reqBody = {
    firstName: "ali",
    lastName: "raza",
    eSimNumber: "# 40547-04-02-1997-5650200",
};

describe("Post Activate Users", () => {
    beforeEachPostCall("/user/active_user", reqBody, nockResponse, 200);
    it("post activate users", async () => {
        const response = await gCall({
            source: POST_ACTIVATE_USER_MUTATION,
            variableValues: {
                input: reqBody,
            },
            contextValue: {
                req: {
                    headers: {
                        authorisation: "test",
                    },
                },
            },
        });
        expect(response).toMatchObject({
            data: {
                activateUser: {
                    success: true,
                },
            },
        });
    });
});
