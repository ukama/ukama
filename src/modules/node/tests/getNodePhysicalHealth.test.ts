import "reflect-metadata";
import { gCall, beforeEachGetCall } from "../../../common/utils";
import { GET_NODE_PHYSICAL_HEALTH_QUERY } from "../../../common/graphql";
import { HEADER } from "../../../constants";

const nockResponse = {
    status: "success",
    data: { temperature: 11, memory: 14, cpu: 12, io: 9 },
};

describe("Get Node physical Health ", () => {
    beforeEachGetCall("/node/physical_health", nockResponse, 200);
    it("get node physical health", async () => {
        const response = await gCall({
            source: GET_NODE_PHYSICAL_HEALTH_QUERY,
            contextValue: {
                req: HEADER,
            },
        });

        expect(response).toMatchObject({
            data: {
                getNodePhysicalHealth: nockResponse.data,
            },
        });
    });
});
