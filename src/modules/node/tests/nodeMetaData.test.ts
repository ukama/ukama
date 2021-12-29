import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_NODE_META_DATA_QUERY } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: { throughput: 11, usersAttached: 1 },
};

describe("Get Node Meta Data ", () => {
    beforeEachGetCall("/node/meta_data", nockResponse, 200);
    it("get node meta data", async () => {
        const response = await gCall({
            source: GET_NODE_META_DATA_QUERY,
            contextValue: {
                req: HEADER,
            },
        });

        expect(response).toMatchObject({
            data: {
                getNodeMetaData: nockResponse.data,
            },
        });
    });
});
