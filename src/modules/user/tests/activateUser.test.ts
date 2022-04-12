import "reflect-metadata";
import { gCall, beforeEachPostCall } from "../../../common/utils";
import { POST_ACTIVATE_USER_MUTATION } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = { status: "success", data: { success: true } };
const reqBody = {
    name: "Ali Raza",
    eSimNumber: "# 40547-04-02-1997-5650200",
    iccid: "# 40547-04-02-1997-5650200",
    roaming: false,
    dataUsage: 1,
    dataPlan: 3,
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
                req: HEADER,
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
