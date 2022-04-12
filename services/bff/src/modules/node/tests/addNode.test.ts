import "reflect-metadata";
import { gCall, beforeEachPostCall } from "../../../common/utils";
import { POST_ADD_NODE_MUTATION } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = { status: "success", data: { success: true } };
const reqBody = {
    name: " Abc Node",
    serialNo: "# 123",
    securityCode: "1234",
};

describe("POST Add Node", () => {
    beforeEachPostCall("/node/add_node", reqBody, nockResponse, 200);
    it("post add node", async () => {
        const response = await gCall({
            source: POST_ADD_NODE_MUTATION,
            variableValues: {
                input: reqBody,
            },
            contextValue: {
                req: HEADER,
            },
        });
        // expect(response).toMatchObject({
        //     data: {
        //         addNode: {
        //             success: true,
        //         },
        //     },
        // });
    });
});
