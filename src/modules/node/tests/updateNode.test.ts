import "reflect-metadata";
import { gCall, beforeEachPostCall } from "../../../common/utils";
import { POST_UPDATE_NODE_MUTATION } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: { id: "1423546576cgfhvgjhb", name: "abc", serialNo: "#4038554" },
};
const reqBody = {
    id: "1423546576cgfhvgjhb",
    name: "abc",
    securityCode: "1234",
};

describe("POST Update Node", () => {
    beforeEachPostCall("/node/update_node", reqBody, nockResponse, 200);
    it("post update node", async () => {
        const response = await gCall({
            source: POST_UPDATE_NODE_MUTATION,
            variableValues: {
                input: reqBody,
            },
            contextValue: {
                req: HEADER,
            },
        });

        // expect(response).toMatchObject({
        //     data: {
        //         updateNode: nockResponse.data,
        //     },
        // });
    });
});
