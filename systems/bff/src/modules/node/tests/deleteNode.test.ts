import "reflect-metadata";
import { gCall, beforeEachPostCall } from "../../../common/utils";
import { DELETE_NODE_MUTATION } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: { id: "ewsrdt11", success: true },
};
const reqBody = {
    id: "ewsrdt11",
};

describe("Delete Node", () => {
    beforeEachPostCall("/node/delete_node", reqBody, nockResponse, 200);
    it("delete node", async () => {
        const response = await gCall({
            source: DELETE_NODE_MUTATION,
            variableValues: {
                input: reqBody.id,
            },
            contextValue: {
                req: HEADER,
            },
        });
        expect(response).toMatchObject({
            data: {
                deleteNode: nockResponse.data,
            },
        });
    });
});
