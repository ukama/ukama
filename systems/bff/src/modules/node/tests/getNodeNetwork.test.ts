import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_NODE_NETWORK } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: {
        id: "55b47a90-7223-44e7-8916-d215110eb1f6",
        status: "ONLINE",
        description: "21 days 5 hours 1 minute",
    },
};

describe("Get Node Network ", () => {
    beforeEachGetCall("/node/get_network", nockResponse, 200);
    it("get node network", async () => {
        const response = await gCall({
            source: GET_NODE_NETWORK,
            contextValue: {
                req: HEADER,
            },
        });

        expect(response).toMatchObject({
            data: {
                getNodeNetwork: nockResponse.data,
            },
        });
    });
});
