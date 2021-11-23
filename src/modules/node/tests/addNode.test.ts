import "reflect-metadata";
import { gCall, beforeEachPostCall } from "../../../common/utils";
import { POST_ADD_NODE_MUTATION } from "../../../common/graphql";

const nockResponse = { status: "success", data: { success: true } };
const reqBody = {
    name: " Abc Node",
    serialNo: "# 123",
};

describe("Post Activate Users", () => {
    beforeEachPostCall("/node/add_node", reqBody, nockResponse, 200);
    it("post activate users", async () => {
        const response = await gCall({
            source: POST_ADD_NODE_MUTATION,
            variableValues: {
                input: reqBody,
            },
            contextValue: {
                req: {
                    headers: {
                        csrf_token: "test",
                        kratos_session: "test",
                    },
                },
            },
        });
        expect(response).toMatchObject({
            data: {
                addNode: {
                    success: true,
                },
            },
        });
    });
});
