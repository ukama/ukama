import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_NODE_DETAILS_QUERY } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "4852428a-6f39-443f-8f43-141fae7a4c00",
        modelType: "Work Node",
        serial: 8605131844265595000,
        macAddress: 9897499773651362000,
        osVersion: 8,
        manufacturing: 8552429119331969,
        ukamaOS: 6,
        hardware: 6,
        description: "Home node is a xyz",
    },
};

describe("Get Node Details ", () => {
    beforeEachGetCall("/node/node_details", nockResponse, 200);
    it("get node details", async () => {
        const response = await gCall({
            source: GET_NODE_DETAILS_QUERY,
            contextValue: {
                req: HEADER,
            },
        });

        expect(response).toMatchObject({
            data: {
                getNodeDetails: nockResponse.data,
            },
        });
    });
});
